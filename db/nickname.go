package db

import (
	basic_errors "errors"
	"fmt"

	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Nickname struct {
	gorm.Model
	DiscordUserId  int64
	DiscordGuildId int64
	Name           string
}

func FindNickname(discordUserId, discordGuildId int64) (*Nickname, error) {
	if r, found := mc.Get(getNicknameCacheKey(discordUserId, discordGuildId)); found {
		return r.(*Nickname), nil
	}

	nickname := Nickname{}
	if err := dbs.Take(&nickname, "discord_user_id=? AND discord_guild_id=?", discordUserId, discordGuildId).Error; err != nil {
		if basic_errors.Is(err, gorm.ErrRecordNotFound) {
			return &nickname, nil
		}
		return nil, errors.WithStack(err)
	}

	mc.Set(
		getNicknameCacheKey(nickname.DiscordUserId, nickname.DiscordGuildId),
		&nickname,
		cache.DefaultExpiration,
	)
	return &nickname, nil
}

func UpdateNickname(discordUserId, discordGuildId int64, name string) (*Nickname, error) {
	nickname, err := FindNickname(discordUserId, discordGuildId)
	if err != nil {
		return nil, err
	}
	if nickname.ID == 0 || nickname.Name != name {
		nickname.DiscordUserId = discordUserId
		nickname.DiscordGuildId = discordGuildId
		nickname.Name = name
		if err := dbs.Save(nickname).Error; err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return nickname, nil
}

func getNicknameCacheKey(discordUserId, discordGuildId int64) string {
	return fmt.Sprintf("nickname/discord_user_id=%v&discord_guild_id=%v", discordUserId, discordGuildId)
}
