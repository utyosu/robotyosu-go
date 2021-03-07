package slack

import (
	"fmt"
	"github.com/k0kubun/pp/v3"
	"github.com/slack-go/slack"
	"log"
)

type Config struct {
	Channel string
	Token   string
	Title   string
}

func (c *Config) Post(p ...interface{}) {
	log.Printf("[%v]\n%v", c.Title, printString(true, p...))
	client := slack.New(c.Token)
	c.postSlack(client, p...)
}

func (c *Config) postSlack(client *slack.Client, p ...interface{}) {
	_, err := client.UploadFile(
		slack.FileUploadParameters{
			Title:    c.Title,
			Content:  printString(false, p...),
			Channels: []string{c.Channel},
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
