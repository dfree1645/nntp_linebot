package model

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"testing"
	"time"
)

func TestPushArticles(t *testing.T) {
	channelsecret := "a"
	channeltoken := "a"
	line, err := linebot.New(channelsecret, channeltoken)
	if err != nil {
		t.Fatal(err)
	}

	articles := []Article{}
	articles = append(articles, Article{ID: 12345, IDstr: "234234@hogehoge", Group: &Group{ID: 1, Name: "testGroup"}, SendDate: time.Now(), Subject: "件名", Body: "本文ほんぶんホンブン"})
	users := []User{}
	users = append(users, User{ID: 123, Name: "表示名", LineID: ""})
	err = PushArticles(line, articles, users, 10)
	if err != nil {
		t.Fatal(err)
	}
}
