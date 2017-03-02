package base

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dfree1645/nntp_linebot/controller"
	"github.com/dfree1645/nntp_linebot/db"
	//"github.com/dfree1645/nntp_linebot/nntp"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/line/line-bot-sdk-go/linebot"
	//"github.com/robfig/cron"
	csrf "github.com/utrack/gin-csrf"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v1"
)

// Serverはベースアプリケーションのserverを示します
type Server struct {
	dbx        *gorm.DB
	sshServer  string
	sshConfig  *ssh.ClientConfig
	nntpServer string
	line       *linebot.Client
	Engine     *gin.Engine
	//cronObj    *cron.Cron
}

func (s *Server) Close() error {
	return s.dbx.Close()
}

// InitはServerを初期化する
func (s *Server) Init(conf, dbconf, env, path string) {
	cs, err := db.NewConfigsFromFile(dbconf)
	if err != nil {
		log.Fatalf("cannot open database configuration. exit. %s", err)
	}
	dbx, err := cs.Open(env)
	if err != nil {
		log.Fatalf("db initialization failed: %s", err)
	}
	s.dbx = dbx

	if err = s.configFromFile(conf); err != nil {
		log.Fatalf("cannot open configuration. exit. %s", err)
	}

	//s.cronObj = cron.New()

	store := sessions.NewCookieStore([]byte("secret"))
	s.Engine.Use(sessions.Sessions("session", store))
	s.Engine.Use(csrf.Middleware(csrf.Options{
		Secret: "secret",
		ErrorFunc: func(c *gin.Context) {
			c.JSON(400, gin.H{"error": "CSRF token mismach"})
			c.Abort()
		},
	}))

	log.Printf("\n%# v\n", s.line)
	s.Route(path)
}

// Newはベースアプリケーションを初期化します
func New(path string) *Server {
	r := gin.Default()
	//r.LoadHTMLGlob(path + "templates/*")
	return &Server{Engine: r}
}

func (s *Server) Run(addr ...string) {
	s.Engine.Run(addr...)
}

func (s *Server) configFromFile(path string) error {
	type Data struct {
		SSHserver   string            `yaml:"sshserver"`
		SSHuser     string            `yaml:"sshuser"`
		SSHpassword string            `yaml:"sshpassword"`
		NNTPserver  string            `yaml:"nntpserver"`
		Line        map[string]string `yaml:"line"`
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	//log.Printf("\n%s\n", b)
	var conf Data
	if err = yaml.Unmarshal(buf, &conf); err != nil {
		return err
	}

	sshConfig := &ssh.ClientConfig{
		User: conf.SSHuser,
		Auth: []ssh.AuthMethod{
			ssh.Password(conf.SSHpassword),
		},
	}

	s.sshServer = conf.SSHserver
	s.sshConfig = sshConfig
	s.nntpServer = conf.NNTPserver
	s.line, err = linebot.New(conf.Line["channelsecret"], conf.Line["channeltoken"])
	if err != nil {
		return err
	}

	return nil
}

// Routeはベースアプリケーションのroutingを設定します
func (s *Server) Route(path string) {
	// ヘルスチェック用
	s.Engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", "pong")
	})
	s.Engine.GET("/token", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"token": csrf.GetToken(c),
		})
	})

	application := &controller.Application{DB: s.dbx}
	article := &controller.Article{DB: s.dbx}
	cron := &controller.Cron{DB: s.dbx, SSHserver: s.sshServer, SSHconfig: s.sshConfig, NNTPserver: s.nntpServer, Line: s.line}
	line := &controller.Line{DB: s.dbx, Line: s.line}

	s.Engine.Static("/static", path+"/public")

	//admin := s.Engine.Group("/admin")
	//admin.Use(controller.AuthRequired())
	{
		//admin.GET("/", application.GetAdminPage)
	}

	app := s.Engine.Group("/")
	{
		app.GET("/", application.RootPage)
		app.GET("/article/:id/:id2", article.ArticlePage)
		app.GET("/cronJob", cron.Job)
		app.POST("/line/webhock", line.Webhook) //URL間違い /webhock -> /webhook
	}

	//s.cronObj.AddFunc("* */10 * * * *", func() {
	/*		log.Println("** cron **")
		cron.CronJob()
	})
	s.cronObj.Start()
	*/
}
