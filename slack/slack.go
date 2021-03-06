package slack

import (
	"fmt"
	"github.com/slack-go/slack"
	"log"
)

var (
	config *SlackConfig
)

type SlackConfig struct {
	WarningChannel string
	AlertChannel   string
	Token          string
	Title          string
}

func Init(c *SlackConfig) {
	config = c
}

func PostSlackWarning(msg interface{}) {
	log.Printf("warning: %+v\n", msg)
	if config == nil {
		return
	}
	client := slack.New(config.Token)
	postSlack(client, config.WarningChannel, msg)
}

func PostSlackAlert(msg interface{}) {
	log.Printf("alert: %+v\n", msg)
	if config == nil {
		return
	}
	client := slack.New(config.Token)
	postSlack(client, config.AlertChannel, msg)
}

func postSlack(client *slack.Client, channel string, msg interface{}) {
	_, err := client.UploadFile(
		slack.FileUploadParameters{
			Title:    config.Title,
			Content:  fmt.Sprintf("%+v", msg),
			Channels: []string{channel},
		},
	)
	if err != nil {
		log.Println(err)
	}
}
