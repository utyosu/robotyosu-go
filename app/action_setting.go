package app

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/utyosu/robotyosu-go/db"
	"github.com/utyosu/robotyosu-go/i18n"
	"strconv"
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

		// 募集期間制限の参照
	case command.match("reserve_limit_time"):
		sendMessage(
			m.ChannelID,
			fmt.Sprintf(
				"Current reserve limit time is %vsec (0 is unlimited)",
				channel.ReserveLimitTime,
			),
		)
		return true, nil

	// 募集期間制限の変更
	case command.match("reserve_limit_time", "*"):
		reserveLimitTimeString := command.fetch(1)
		reserveLimitTime, err := strconv.ParseUint(reserveLimitTimeString, 10, 32)
		if err != nil {
			sendMessage(m.ChannelID, fmt.Sprintf("Invalid time: %v", reserveLimitTimeString))
			return true, nil
		}
		if err := channel.UpdateChannelReserveLimitTime(uint32(reserveLimitTime)); err != nil {
			return true, err
		}
		sendMessage(
			m.ChannelID,
			fmt.Sprintf(
				"Reserve limit time changed to %vsec (0 is unlimited)",
				reserveLimitTimeString,
			),
		)
		return true, nil

	// 通常の募集期限の参照
	case command.match("expire_duration"):
		sendMessage(
			m.ChannelID,
			fmt.Sprintf(
				"Current expire duration is %vsec",
				channel.ExpireDuration,
			),
		)
		return true, nil

	// 通常の募集期限の変更
	case command.match("expire_duration", "*"):
		expireDurationString := command.fetch(1)
		expireDuration, err := strconv.ParseUint(expireDurationString, 10, 32)
		if err != nil {
			sendMessage(m.ChannelID, fmt.Sprintf("Invalid time: %v", expireDurationString))
			return true, nil
		}
		if err := channel.UpdateExpireDuration(uint32(expireDuration)); err != nil {
			return true, err
		}
		sendMessage(
			m.ChannelID,
			fmt.Sprintf(
				"Expire duration changed to %vsec",
				expireDurationString,
			),
		)
		return true, nil

	// 予約の募集期限の参照
	case command.match("expire_duration_for_reserve"):
		sendMessage(
			m.ChannelID,
			fmt.Sprintf(
				"Current expire duration for reserve is %vsec",
				channel.ExpireDurationForReserve,
			),
		)
		return true, nil

	// 予約の募集期限の変更
	case command.match("expire_duration_for_reserve", "*"):
		expireDurationForReserveString := command.fetch(1)
		expireDurationForReserve, err := strconv.ParseUint(expireDurationForReserveString, 10, 32)
		if err != nil {
			sendMessage(m.ChannelID, fmt.Sprintf("Invalid time: %v", expireDurationForReserveString))
			return true, nil
		}
		if err := channel.UpdateExpireDurationForReserve(uint32(expireDurationForReserve)); err != nil {
			return true, err
		}
		sendMessage(
			m.ChannelID,
			fmt.Sprintf(
				"Expire duration for reserve changed to %vsec",
				expireDurationForReserveString,
			),
		)
		return true, nil

	// 募集機能のヘルプ
	case command.match("help"):
		sendMessage(m.ChannelID, i18n.HelpRecruitmentCommands(i18n.ToLanguage(channel.Language)))
		return true, nil
	}

	return false, nil
}
