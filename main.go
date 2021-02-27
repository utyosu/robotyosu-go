package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/utyosu/robotyosu-go/db"
	"github.com/utyosu/robotyosu-go/env"
	"github.com/utyosu/robotyosu-go/i18n"
	"log"
	"strconv"
	"time"
)

const (
	commandPrefix = ".rt"
)

var (
	discordSession *discordgo.Session
	stopBot        = make(chan bool)
)

func init() {
	if loc, err := time.LoadLocation("Asia/Tokyo"); err == nil {
		time.Local = loc
	}
}

func main() {
	defer notifySlackWhenPanic("main")

	dbs := db.ConnectDb()
	defer dbs.Close()

	var err error
	discordSession, err = discordgo.New()
	if err != nil {
		panic(err)
	}
	discordSession.Token = fmt.Sprintf("Bot %s", env.DiscordBotToken)

	discordSession.AddHandler(onMessageCreate)
	if err = discordSession.Open(); err != nil {
		panic(err)
	}
	defer discordSession.Close()

	log.Println("Listening...")

	doFuncSchedule(closeExpiredRecruitment, time.Second*env.ScheduledDuration)
	doFuncSchedule(notifyReservedRecruitmentOnTime, time.Second*env.ScheduledDuration)

	<-stopBot
	return
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer notifySlackWhenPanic(messageInformation(s, m))

	// 自分のメッセージは処理しない
	if m.Author.ID == env.DiscordBotClientId {
		return
	}
	log.Printf("\t%v\t%v\t%v\t%v\t%v\n", m.GuildID, m.ChannelID, m.Type, m.Author.Username, m.Content)

	// サーバーIDがない(=DM)は処理しない
	if m.GuildID == "" {
		return
	}

	discordChannelId, _ := strconv.ParseInt(m.ChannelID, 10, 64)

	// 全チャンネルで使えるコマンド
	if err := actionAllChannel(s, m, discordChannelId); err != nil {
		sendMessage(m.ChannelID, i18n.T(i18n.DefaultLanguage, "error"))
		postSlackWarning(err)
		return
	}

	channel, err := db.FindChannel(discordChannelId)
	if err != nil {
		postSlackWarning(err)
		return
	}

	if !channel.IsValidAnyFunction() {
		return
	}

	// 何かしらの機能が有効なチャンネルで使えるコマンド
	if err := actionValidChannel(s, m, channel); err != nil {
		postSlackWarning(err)
		sendMessageT(channel, "error")
		return
	}

	authorId, _ := strconv.ParseInt(m.Author.ID, 10, 64)
	userName := m.Author.Username
	if m.Member != nil && m.Member.Nick != "" {
		userName = m.Member.Nick
	}
	user, err := db.FindOrCreateUser(authorId, userName)
	if err != nil {
		sendMessageT(channel, "error")
		postSlackWarning(err)
		return
	}

	// recruitmentが有効なチャンネルで使えるコマンド
	if channel.Recruitment {
		if err := actionRecruitment(s, m, channel, user); err != nil {
			sendMessageT(channel, "error")
			postSlackWarning(err)
			return
		}
	}
}

func sendMessageT(c *db.Channel, key string, params ...interface{}) {
	sendMessage(c.DiscordIdStr(), i18n.T(c.Language, key, params...))
}

func sendMessage(channelID string, msg string) {
	if _, err := discordSession.ChannelMessageSend(channelID, msg); err != nil {
		postSlackWarning(fmt.Sprintf("Error sending message: %v", err))
	}
}
