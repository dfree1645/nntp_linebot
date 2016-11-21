package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	//"github.com/dfree1645/nntp_linebot/model"
)

type Article struct {
	DB *gorm.DB
}

func (a *Article) ArticlePage(c *gin.Context) {
	// id : DB上のオートインクリメント番号
	// id2: NNTPサーバー上のID

	id1, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "wrong article ids\n")
	}
	id2 := c.Param("id2")
	c.String(http.StatusOK, "This is articlePage.\n\nid1=%d\nid2=%s\n", id1, id2)
}
