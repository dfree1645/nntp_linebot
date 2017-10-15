package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/line/line-bot-sdk-go/linebot"
	//"github.com/dfree1645/nntp_linebot/model"
)

type Line struct {
	DB   *gorm.DB
	Line *linebot.Client
}

func (a *Line) Webhook(c *gin.Context) {
	received, err := a.Line.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			log.Println(err)
		}
		return
	}

	for _, event := range received {
		log.Println()
		log.Printf("message received\n UserID: %s\n type: %s\n\n", event.Source.UserID, event.Type)
		if event.Type == linebot.EventTypeFollow {
			message := linebot.NewTextMessage("FollowEvent\nYour UserID: " + event.Source.UserID)
			_, err := a.Line.ReplyMessage(event.ReplyToken, message).Do()
			if err != nil {
				log.Println(err)
			}
			log.Println("Follow userID:", event.Source.UserID)
		} else if event.Type == linebot.EventTypeMessage {
			message := linebot.NewTextMessage("MessageEvent\nYour UserID: " + event.Source.UserID)
			_, err := a.Line.ReplyMessage(event.ReplyToken, message).Do()
			if err != nil {
				log.Println(err)
			}
		} else if event.Type == linebot.EventTypeJoin {
			message := linebot.NewTextMessage("JoinEvent\nGroup UserID: " + event.Source.GroupID + "\nRoomID: " + event.Source.RoomID)
			_, err := a.Line.ReplyMessage(event.ReplyToken, message).Do()
			if err != nil {
				log.Println(err)
			}
		}
	}
}
