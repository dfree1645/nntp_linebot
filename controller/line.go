package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/line/line-bot-sdk-go/linebot"
	//"github.com/dfree1645/nntp_linebot/model"
)

type Line struct {
	DB   *gorm.DB
	Line *linebot.Client
}

func (a *Line) Webhock(c *gin.Context) {
	c.String(http.StatusOK, "This is LineWebhock\n")
}
