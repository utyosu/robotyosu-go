package db

import (
	basic_errors "errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	DiscordUserId int64
	Name          string
	Nickname      Nickname `gorm:"foreignkey:DiscordUserId;references:discord_user_id"`
}

func FindUser(discordUserId int64) (*User, error) {
	if r, found := mc.Get(getUserCacheKey(discordUserId)); found {
		return r.(*User), nil
	}

	user := User{}
	if err := dbs.Take(&user, "discord_user_id=?", discordUserId).Error; err != nil {
		if !basic_errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(err)
		}
	}

	mc.Set(
		getUserCacheKey(user.DiscordUserId),
		&user,
		cache.DefaultExpiration,
	)
	return &user, nil
}

func FindOrCreateUser(discordUserId int64, name string) (*User, error) {
	user, err := FindUser(discordUserId)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		user.DiscordUserId = discordUserId
		user.Name = name
		if err := dbs.Create(&user).Error; err != nil {
			return nil, errors.WithStack(err)
		}
	} else if user.Name != name {
		user.Name = name
		if err := dbs.Save(&user).Error; err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return user, nil
}

func (u *User) DisplayName() string {
	if u.Nickname.Name != "" {
		return u.Nickname.Name
	}
	return u.Name
}

func getUserCacheKey(discordUserId int64) string {
	return fmt.Sprintf("user/discord_user_id=%v", discordUserId)
}
