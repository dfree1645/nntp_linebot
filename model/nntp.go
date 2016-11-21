package model

import (
	"github.com/dfree1645/nntp_linebot/nntp"
	"github.com/jinzhu/gorm"
	"time"
)

type Group struct {
	ID   int64 `gorm:"primary_key"`
	Name string
	High int64
	Low  int64
}

type Article struct {
	ID       int64
	IDstr    string
	Group    *Group
	SendDate time.Time
	Subject  string
	Body     string
}

type User struct {
	ID     int64 `sql:"user_id"`
	Name   string
	LineID string
}

func GetNewsGroups(db *gorm.DB) ([]Group, error) {
	groups := []Group{}
	d := db.Table("groups").Find(&groups)
	if d.Error != nil {
		return []Group{}, d.Error
	}
	return groups, nil
}

func ConvToArticle(org *nntp.Article, group *Group) Article {
	// TODO:文字コードの処理など
	return Article{ID: 12345, IDstr: "234234@hogehoge", Group: group, SendDate: time.Now(), Subject: "件名", Body: "本文ほんぶんホンブン"}
}

func GetUsers(db *gorm.DB, group *Group) ([]User, error) {
	// TODO:購読するグループを選択できるようにする
	// 現在は全有効ユーザ一覧を返す
	users := []User{}
	d := db.Table("users").Find(&users)
	if d.Error != nil {
		return nil, d.Error
	}
	return users, nil
}

func UpdateGroup(db *gorm.DB, group *Group, high, low int) error {
	group.High = int64(high)
	group.Low = int64(low)
	db.Table("groups").Update(group)
	return nil
}

func InsertArticles(db *gorm.DB, articles []Article) error {
	// TODO:実装
	return nil
}
