package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/dfree1645/nntp_linebot/model"
	"github.com/dfree1645/nntp_linebot/nntp"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/crypto/ssh"
)

type Cron struct {
	DB         *gorm.DB
	SSHserver  string
	SSHconfig  *ssh.ClientConfig
	NNTPserver string
	Line       *linebot.Client
}

func (a *Cron) Job(c *gin.Context) {
	err := a.CronJob()
	if err != nil {
		c.String(500, err.Error())
	}
	c.String(http.StatusOK, "This is CronJob\n")
}

func (a *Cron) CronJob() error {
	log.Println("-- CronJob start. --")
	defer log.Println("-- CronJob finish. --")

	groups, err := model.GetNewsGroups(a.DB)
	if err != nil {
		log.Printf("error: %s\n", err.Error())
		return err
	}
	log.Printf("groups: %s\n", len(groups))
	if len(groups) == 0 {
		return nil
	}

	client, err := ssh.Dial("tcp", a.SSHserver, a.SSHconfig)
	if err != nil {
		log.Println("error SSH: " + err.Error())
		return err
	}
	defer client.Close()

	nsC, err := client.Dial("tcp", a.NNTPserver)
	if err != nil {
		log.Println("error NNTPconn: " + err.Error())
		return err
	}
	defer nsC.Close()

	nsNntp, err := nntp.New(nsC)
	if err != nil {
		log.Println("error NNTPsess: " + err.Error())
		return err
	}

	for _, v := range groups {
		// connect to a news group
		_, low, high, err := nsNntp.Group(v.Name)
		if err != nil {
			log.Printf("[%s] Could not connect to groups: %v\n", v.Name, err)
			continue
		}

		// 新着ニュース有無
		if v.High == int64(high) {
			log.Printf("[%s] There are no new articles.\n", v.Name)
			continue
		}
		// 新着ニュースのID範囲
		newHigh := int64(high)
		newLow := v.High + 1
		if newLow < int64(low) {
			newLow = int64(low)
		}
		newArticles := []model.Article{}
		for i := newHigh; i >= newLow; i-- {
			article, err := nsNntp.Article(strconv.FormatInt(i, 10))
			if err != nil {
				log.Printf("[%s] Could not fetch article (id=%d) : %v\n", v.Name, i, err)
				continue
			}
			newArticles = append(newArticles, model.ConvToArticle(article, &v))
		}
		log.Printf("[%s] %d new articles\n", v.Name, len(newArticles))

		// このグループ購読中ユーザー一覧取得
		users, err := model.GetUsers(a.DB, &v)
		if err != nil {
			log.Printf("[%s] Error! Could not get users: %s\n", v.Name, err.Error())
		}

		if err := model.PushArticles(a.Line, newArticles, users, 10); err != nil {
			log.Printf("[%s] Error! Could not push line message: %s\n", v.Name, err.Error())
		}

		if err := model.InsertArticles(a.DB, newArticles); err != nil {
			log.Printf("[%s] Error! Could not insert articles into DB: %s\n", v.Name, err.Error())
		}

		if err := model.UpdateGroup(a.DB, &v, high, low); err != nil {
			log.Printf("[%s] Error! Could not update DB: %s\n", v.Name, err.Error())
		}

	}
	return nil
}
