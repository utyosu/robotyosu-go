package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/utyosu/robotyosu-go/db"
	"github.com/utyosu/robotyosu-go/i18n"
)

func actionAllChannel(s *discordgo.Session, m *discordgo.MessageCreate, discordChannelId int64) error {
	switch {
	// 有効化
	case m.Content == (commandPrefix + " enable"):
		channel, err := db.FindOrCreateChannel(discordChannelId)
		if err != nil {
			return err
		}
		if err := channel.UpdateChannelRecruitment(true); err != nil {
			return err
		}
		sendMessage(m.ChannelID, i18n.CommonMessage("enable"))

	// 無効化
	case m.Content == (commandPrefix + " disable"):
		channel, err := db.FindOrCreateChannel(discordChannelId)
		if err != nil {
			return err
		}
		if err := channel.UpdateChannelRecruitment(false); err != nil {
			return err
		}
		sendMessage(m.ChannelID, i18n.CommonMessage("disable"))

	// コマンドヘルプの表示
	case m.Content == (commandPrefix + " help"):
		language := i18n.DefaultLanguage

		// チャンネルが見つかれば言語を設定する
		if channel, err := db.FindChannel(discordChannelId); err != nil {
			return err
		} else if channel.ID != 0 {
			language = channel.Language
		}
		sendMessage(m.ChannelID, i18n.HelpBasicCommands(language))
	}
	return nil
}
