package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Application struct {
	DB *gorm.DB
}

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func (a *Application) RootPage(c *gin.Context) {
	c.String(http.StatusOK, "Hello, world!\n")
}
