package app

import (
	"fmt"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/utyosu/robotyosu-go/db"
	"github.com/utyosu/robotyosu-go/env"
	"github.com/utyosu/robotyosu-go/i18n"
	"github.com/utyosu/robotyosu-go/slack"
)

var (
	discordSession *discordgo.Session
	stopBot        = make(chan bool)
	slackAlert     *slack.Config
	slackWarning   *slack.Config
	commitHash     string
	buildDatetime  string
)

func init() {
	slackAlert = &slack.Config{
		Channel: env.SlackChannelAlert,
		Token:   env.SlackToken,
		Title:   env.SlackTitleAlert,
	}
	slackWarning = &slack.Config{
		Channel: env.SlackChannelWarning,
		Token:   env.SlackToken,
		Title:   env.SlackTitleWarning,
	}
}

func Start() {
	var err error
	discordSession, err = discordgo.New(fmt.Sprintf("Bot %s", env.DiscordBotToken))
	if err != nil {
		panic(err)
	}

	discordSession.AddHandler(onMessageCreate)
	if err = discordSession.Open(); err != nil {
		panic(err)
	}
	defer discordSession.Close()
	log.Println("Listening...")

	doFuncSchedule(closeExpiredRecruitment, env.ScheduledDuration)
	doFuncSchedule(notifyReservedRecruitmentOnTime, env.ScheduledDuration)
	<-stopBot
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer NotifySlackWhenPanic(s, m)

	// 自分のメッセージは処理しない
	if m.Author.ID == env.DiscordBotClientId {
		return
	}
	log.Printf("\t%v\t%v\t%v\t%v\t%v\n", m.GuildID, m.ChannelID, m.Type, m.Author.Username, m.Content)

	// サーバーIDがない(=DM)は処理しない
	if m.GuildID == "" {
		return
	}

	// 全チャンネルで使えるコマンド
	if processed, err := actionAllChannel(s, m); err != nil {
		sendMessage(m.ChannelID, i18n.T(i18n.DefaultLanguage, "error"))
		slackWarning.Post(err)
		return
	} else if processed {
		return
	}

	discordChannelId, _ := strconv.ParseInt(m.ChannelID, 10, 64)
	channel, err := db.FindChannel(discordChannelId)
	if err != nil {
		slackWarning.Post(err)
		return
	}

	if !channel.IsEnabledRecruitment() {
		return
	}

	// 何かしらの機能が有効なチャンネルで使えるコマンド
	if processed, err := actionSetting(s, m, channel); err != nil {
		slackWarning.Post(err)
		sendMessageT(channel, "error")
		return
	} else if processed {
		return
	}

	authorId, _ := strconv.ParseInt(m.Author.ID, 10, 64)
	userName := m.Author.Username
	user, err := db.FindOrCreateUser(authorId, userName)
	if err != nil {
		sendMessageT(channel, "error")
		slackWarning.Post(err)
		return
	}
	discordGuildId, _ := strconv.ParseInt(m.GuildID, 10, 64)
	nickname, err := db.UpdateNickname(user.DiscordUserId, discordGuildId, getNickname(m))
	if err != nil {
		sendMessageT(channel, "error")
		slackWarning.Post(err)
		return
	}
	user.Nickname = *nickname

	// recruitmentが有効なチャンネルで使えるコマンド
	if channel.Recruitment {
		if err := actionRecruitment(s, m, channel, user); err != nil {
			sendMessageT(channel, "error")
			slackWarning.Post(err)
			return
		}
	}
}

func sendMessageT(c *db.Channel, key string, params ...interface{}) {
	sendMessage(c.DiscordIdStr(), i18n.T(c.Language, key, params...))
}

func sendMessage(channelID string, msg string) {
	if _, err := discordSession.ChannelMessageSend(channelID, msg); err != nil {
		slackWarning.Post(
			errors.New("Error sending message"),
			msg,
			channelID,
		)
	}
}

func getNickname(m *discordgo.MessageCreate) string {
	if m.Member != nil && m.Member.Nick != "" {
		return m.Member.Nick
	}
	if m.Author.GlobalName != "" {
		return m.Author.GlobalName
	}
	return m.Author.Username
}
