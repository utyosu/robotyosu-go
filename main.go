package main

import (
	"github.com/utyosu/robotyosu-go/app"
	"github.com/utyosu/robotyosu-go/db"
	"github.com/utyosu/robotyosu-go/env"
	"github.com/utyosu/robotyosu-go/slack"
	"time"
)

func init() {
	if loc, err := time.LoadLocation("Asia/Tokyo"); err == nil {
		time.Local = loc
	}
	slack.Init(&slack.SlackConfig{
		WarningChannel: env.SlackChannelWarning,
		AlertChannel:   env.SlackChannelAlert,
		Token:          env.SlackToken,
		Title:          env.SlackTitle,
	})
}

func main() {
	defer app.NotifySlackWhenPanic("main")
	db.ConnectDb()
	app.Start()
	return
}
