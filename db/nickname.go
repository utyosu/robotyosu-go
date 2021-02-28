package db

import (
	basic_errors "errors"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Nickname struct {
	gorm.Model
	UserId         uint
	DiscordGuildId int64
	Name           string
}

func FindNickname(userId uint, guildId int64) (*Nickname, error) {
	nickname := Nickname{}
	if err := dbs.Take(&nickname, "user_id=? AND discord_guild_id=?", userId, guildId).Error; err != nil {
		if basic_errors.Is(err, gorm.ErrRecordNotFound) {
			return &nickname, nil
		}
		return nil, errors.WithStack(err)
	}
	return &nickname, nil
}

func UpdateNickname(userId uint, guildId int64, name string) (*Nickname, error) {
	nickname, err := FindNickname(userId, guildId)
	if err != nil {
		return nil, err
	}
	if nickname.ID == 0 || nickname.Name != name {
		nickname.UserId = userId
		nickname.DiscordGuildId = guildId
		nickname.Name = name
		if err := dbs.Save(nickname).Error; err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return nickname, nil
}
