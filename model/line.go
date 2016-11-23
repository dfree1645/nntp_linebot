package model

import (
	"github.com/line/line-bot-sdk-go/linebot"
)

func PushArticles(line *linebot.Client, articles []Article, users []User, max int) error {
	// 配列articlesを、配列usersのlineIDに配信。ただし、最大件数はmaxとする
	for _, u := range users {
		for _, a := range articles {
			if _, err := line.PushMessage(u.LineID, linebot.NewTextMessage("【Group】\n"+a.Group.Name+"\n【Subject】\n"+a.Subject+"\n【Body】\n"+a.Body+"\nDate:"+a.SendDate.String())).Do(); err != nil {
				return err
			}
		}
	}
	return nil
}
