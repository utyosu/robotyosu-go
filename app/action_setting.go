package app

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/utyosu/robotyosu-go/db"
	"github.com/utyosu/robotyosu-go/i18n"
	"regexp"
	"strings"
	"time"
)

var (
	regexpShowTimezone        = regexp.MustCompile(`\A` + commandPrefix + `\s+timezone\z`)
	regexpSetTimezone         = regexp.MustCompile(`\A` + commandPrefix + `\s+timezone\s+([\w/]+)\z`)
	regexpShowLanguage        = regexp.MustCompile(`\A` + commandPrefix + `\s+language\z`)
	regexpSetLanguage         = regexp.MustCompile(`\A` + commandPrefix + `\s+language\s+([\w/]+)\z`)
	regexpShowRecruitmentHelp = regexp.MustCompile(`使い方|ヘルプ|help`)
)

func actionSetting(s *discordgo.Session, m *discordgo.MessageCreate, channel *db.Channel) (bool, error) {
	switch {
	// タイムゾーンの参照
	case regexpShowTimezone.MatchString(m.Content):
		sendMessage(
			m.ChannelID,
			fmt.Sprintf(
				"Current timezone is %v\nAvailable timezones: UTC, EST, GMT, Asia/Tokyo, etc.\n",
				channel.Timezone,
			),
		)
		return true, nil

	// タイムゾーンの変更
	case regexpSetTimezone.MatchString(m.Content):
		timezoneString := getMatchRegexpString(m.Content, regexpSetTimezone)
		_, err := time.LoadLocation(timezoneString)
		if err != nil {
			sendMessage(m.ChannelID, fmt.Sprintf("No such timezone: %v", timezoneString))
			return true, nil
		}
		if err = channel.UpdateChannelTimezone(timezoneString); err != nil {
			return true, err
		}
		sendMessage(m.ChannelID, fmt.Sprintf("Timezone changed to %v", timezoneString))

	// 言語の参照
	case regexpShowLanguage.MatchString(m.Content):
		sendMessage(
			m.ChannelID,
			fmt.Sprintf(
				"Current language is %v\nAvailable languages: %v",
				i18n.ToLanguage(channel.Language),
				strings.Join(i18n.Languages, ", "),
			),
		)
		return true, nil

	// 言語の変更
	case regexpSetLanguage.MatchString(m.Content):
		languageString := getMatchRegexpString(m.Content, regexpSetLanguage)
		if languageString != i18n.ToLanguage(languageString) {
			sendMessage(m.ChannelID, fmt.Sprintf("No such language: %v", languageString))
			return true, nil
		}
		if err := channel.UpdateChannelLanguage(languageString); err != nil {
			return true, err
		}
		sendMessage(m.ChannelID, fmt.Sprintf("Language changed to %v", languageString))

	// 募集機能のヘルプ
	case regexpShowRecruitmentHelp.MatchString(m.Content):
		sendMessage(m.ChannelID, i18n.HelpRecruitmentCommands(i18n.ToLanguage(channel.Language)))
		return true, nil
	}

	return false, nil
}
