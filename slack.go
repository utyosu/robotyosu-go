package main

import (
	"fmt"
	"github.com/slack-go/slack"
	"github.com/utyosu/robotyosu-go/env"
	"log"
)

func postSlackWarning(msg interface{}) {
	log.Printf("warning: %+v\n", msg)
	if env.SlackToken == "" {
		return
	}
	client := slack.New(env.SlackToken)
	channel := env.SlackChannelWarning
	postSlack(client, channel, msg)
}

func postSlackAlert(msg interface{}) {
	log.Printf("alert: %+v\n", msg)
	if env.SlackToken == "" {
		return
	}
	client := slack.New(env.SlackToken)
	channel := env.SlackChannelAlert
	postSlack(client, channel, msg)
}

func postSlack(client *slack.Client, channel string, msg interface{}) {
	content := fmt.Sprintf("robotyosu-go error notification\n```\n%+v\n```", msg)
	_, _, err := client.PostMessage(
		channel,
		slack.MsgOptionText(content, true),
	)
	if err != nil {
		log.Println(err)
	}
}
