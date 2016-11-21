package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/crypto/ssh"
	//"github.com/dfree1645/nntp_linebot/model"
)

type Cron struct {
	DB         *gorm.DB
	SSHserver  string
	SSHconfig  *ssh.ClientConfig
	NNTPserver string
	Line       *linebot.Client
}

func (a *Cron) Job(c *gin.Context) {
	c.String(http.StatusOK, "This is CronJob\n")
}
