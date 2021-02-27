package db

import (
	basic_errors "errors"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type User struct {
	gorm.Model
	DiscordUserId int64
	Name          string
}

func FindOrCreateUser(discordUserId int64, name string) (*User, error) {
	user := User{}
	if err := dbs.Take(&user, "discord_user_id=?", discordUserId).Error; err != nil {
		// レコードが見つからない以外のエラーならエラーを返す
		if !basic_errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
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
	return &user, nil
}
