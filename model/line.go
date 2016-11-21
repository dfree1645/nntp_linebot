package model

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
)

func PushArticles(line *linebot.Client, articles []Article, users []User, max int) error {
	// 配列articlesを、配列usersのlineIDに配信。ただし、最大件数はmaxとする
	log.Println(articles)
	log.Println(users)
	return nil
}
