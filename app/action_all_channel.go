package app

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/utyosu/robotyosu-go/db"
	"github.com/utyosu/robotyosu-go/i18n"
	"strconv"
)

func actionAllChannel(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	command, ok := toCommand(m.Content)
	if !ok {
		return false, nil
	}

	switch {
	// 有効化
	case command.match("enable"):
		discordChannelId, _ := strconv.ParseInt(m.ChannelID, 10, 64)
		discordGuildId, _ := strconv.ParseInt(m.GuildID, 10, 64)
		channel, err := db.FindOrCreateChannel(discordChannelId, discordGuildId)
		if err != nil {
			return true, err
		}
		if err := channel.UpdateRecruitment(true); err != nil {
			return true, err
		}
		sendMessage(m.ChannelID, i18n.CommonMessage("enable"))
		return true, nil

	// 無効化
	case command.match("disable"):
		discordChannelId, _ := strconv.ParseInt(m.ChannelID, 10, 64)
		discordGuildId, _ := strconv.ParseInt(m.GuildID, 10, 64)
		channel, err := db.FindOrCreateChannel(discordChannelId, discordGuildId)
		if err != nil {
			return true, err
		}
		if err := channel.UpdateRecruitment(false); err != nil {
			return true, err
		}
		sendMessage(m.ChannelID, i18n.CommonMessage("disable"))
		return true, nil

	// コマンドヘルプの表示
	case command.match("help"):
		discordChannelId, _ := strconv.ParseInt(m.ChannelID, 10, 64)
		language := i18n.DefaultLanguage

		// チャンネルが見つかれば言語を設定する
		if channel, err := db.FindChannel(discordChannelId); err != nil {
			return true, err
		} else if channel.ID != 0 {
			language = channel.Language
		}
		sendMessage(m.ChannelID, i18n.HelpBasicCommands(language))
		return true, nil

	// バージョンの表示
	case command.match("version"):
		sendMessage(m.ChannelID, fmt.Sprintf(
			"CommitHash: %v\nBuildDatetime: %v",
			commitHash,
			buildDatetime,
		))
		return true, nil
	}
	return false, nil
}
