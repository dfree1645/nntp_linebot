package base

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/dfree1645/nntp_linebot/controller"
	"github.com/dfree1645/nntp_linebot/db"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/robfig/cron"
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

	// 設定
	log.Printf("SSH Server: %s\n", s.sshServer)
	log.Printf("NNTP Server: %s\n", s.nntpServer)

	s.Route(path)
}

// Newはベースアプリケーションを初期化します
func New(path string) *Server {
	r := gin.Default()
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
		SSHciphers  []string          `yaml:"sshciphers"`
		NNTPserver  string            `yaml:"nntpserver"`
		Line        map[string]string `yaml:"line"`
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var conf Data
	if err = yaml.Unmarshal(buf, &conf); err != nil {
		return err
	}

	sshConfig := &ssh.ClientConfig{
		User: conf.SSHuser,
		Auth: []ssh.AuthMethod{
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
				answers = []string{conf.SSHpassword}
				return
			}),
			ssh.Password(conf.SSHpassword),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	sshConfig.Ciphers = conf.SSHciphers

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

	application := &controller.Application{DB: s.dbx}
	article := &controller.Article{DB: s.dbx}
	croncon := &controller.Cron{DB: s.dbx, SSHserver: s.sshServer, SSHconfig: s.sshConfig, NNTPserver: s.nntpServer, Line: s.line}
	line := &controller.Line{DB: s.dbx, Line: s.line}

	s.Engine.Static("/static", path+"/public")

	app := s.Engine.Group("/")
	{
		app.GET("/", application.RootPage)
		app.GET("/article/:id/:id2", article.ArticlePage)
		app.GET("/cronJob", croncon.Job)
		app.POST("/line/webhook", line.Webhook)
	}

	c := cron.New()
	c.AddFunc("0 */10 * * * *", func() {
		croncon.CronJob()
	})
	c.Start()
}
