package slack

import (
	"fmt"
	"github.com/k0kubun/pp/v3"
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

func PostSlackWarning(p ...interface{}) {
	log.Printf("[warning]\n%v", printString(true, p))
	if config == nil {
		return
	}
	client := slack.New(config.Token)
	postSlack(client, config.WarningChannel, p...)
}

func PostSlackAlert(p ...interface{}) {
	log.Printf("[warning]\n%+v", printString(true, p))
	if config == nil {
		return
	}
	client := slack.New(config.Token)
	postSlack(client, config.AlertChannel, p...)
}

func postSlack(client *slack.Client, channel string, p ...interface{}) {
	_, err := client.UploadFile(
		slack.FileUploadParameters{
			Title:    config.Title,
			Content:  printString(false, p...),
			Channels: []string{channel},
		},
	)
	if err != nil {
		log.Println(err)
	}
}

func printString(enableColor bool, p ...interface{}) string {
	lpp := pp.New()
	lpp.SetColoringEnabled(enableColor)
	var ret string
	for _, v := range p {
		if _, ok := v.(error); ok {
			// errorはppで展開すると読めないのでfmtで表示する
			ret += fmt.Sprintf("%+v\n", v)
		} else {
			ret += lpp.Sprintf("%v\n", v)
		}
	}
	return ret
}
