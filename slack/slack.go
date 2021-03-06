package slack

import (
	"fmt"
	"github.com/slack-go/slack"
	"github.com/utyosu/robotyosu-go/env"
	"log"
)

func PostSlackWarning(msg interface{}) {
	log.Printf("warning: %+v\n", msg)
	if env.SlackToken == "" {
		return
	}
	client := slack.New(env.SlackToken)
	channel := env.SlackChannelWarning
	postSlack(client, channel, msg)
}

func PostSlackAlert(msg interface{}) {
	log.Printf("alert: %+v\n", msg)
	if env.SlackToken == "" {
		return
	}
	client := slack.New(env.SlackToken)
	channel := env.SlackChannelAlert
	postSlack(client, channel, msg)
}

func postSlack(client *slack.Client, channel string, msg interface{}) {
	_, err := client.UploadFile(
		slack.FileUploadParameters{
			Title:    "robotyosu-go error notification",
			Content:  fmt.Sprintf("%+v", msg),
			Channels: []string{channel},
		},
	)
	if err != nil {
		log.Println(err)
	}
}
