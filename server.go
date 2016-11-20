package base

import (
	"io/ioutil"
	"log"
	//"net/http"

	//"github.com/dfree1645/nntp_linebot/controller"
	"github.com/dfree1645/nntp_linebot/db"
	//"github.com/dfree1645/nntp_linebot/nntp"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
	lineConfig map[string]string
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

	store := sessions.NewCookieStore([]byte("secret"))
	s.Engine.Use(sessions.Sessions("session", store))
	s.Engine.Use(csrf.Middleware(csrf.Options{
		Secret: "secret",
		ErrorFunc: func(c *gin.Context) {
			c.JSON(400, gin.H{"error": "CSRF token mismach"})
			c.Abort()
		},
	}))

	log.Printf("\n%# v\n", s)
	log.Printf("\n%s\n", s.lineConfig["channelsecret"])
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
	s.lineConfig = conf.Line

	return nil
}

// Routeはベースアプリケーションのroutingを設定します
func (s *Server) Route(path string) {

}
