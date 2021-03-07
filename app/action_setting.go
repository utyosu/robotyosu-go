package app

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/utyosu/robotyosu-go/db"
	"github.com/utyosu/robotyosu-go/i18n"
	"strings"
	"time"
)

func actionSetting(s *discordgo.Session, m *discordgo.MessageCreate, channel *db.Channel) (bool, error) {
	command, ok := toCommand(m.Content)
	if !ok {
		return false, nil
	}

	switch {
	// タイムゾーンの参照
	case command.match("timezone"):
		sendMessage(
			m.ChannelID,
			fmt.Sprintf(
				"Current timezone is %v\nAvailable timezones: UTC, EST, GMT, Asia/Tokyo, etc.\n",
				channel.Timezone,
			),
		)
		return true, nil

	// タイムゾーンの変更
	case command.match("timezone", "*"):
		timezoneString := command.fetch(1)
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
	case command.match("language"):
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
	case command.match("language", "*"):
		languageString := command.fetch(1)
		if languageString != i18n.ToLanguage(languageString) {
			sendMessage(m.ChannelID, fmt.Sprintf("No such language: %v", languageString))
			return true, nil
		}
		if err := channel.UpdateChannelLanguage(languageString); err != nil {
			return true, err
		}
		sendMessage(m.ChannelID, fmt.Sprintf("Language changed to %v", languageString))

	// 募集機能のヘルプ
	case command.match("help"):
		sendMessage(m.ChannelID, i18n.HelpRecruitmentCommands(i18n.ToLanguage(channel.Language)))
		return true, nil
	}

	return false, nil
}
