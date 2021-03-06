package db

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Participant struct {
	gorm.Model
	DiscordUserId int64
	RecruitmentId uint
	User          User `gorm:"foreignkey:DiscordUserId;references:discord_user_id"`
}

func InsertParticipant(user *User, recruitment *Recruitment) error {
	err := dbs.Create(&Participant{
		DiscordUserId: user.DiscordUserId,
		RecruitmentId: recruitment.ID,
	}).Error
	return errors.WithStack(err)
}

func (p *Participant) Delete() error {
	err := dbs.Delete(p).Error
	return errors.WithStack(err)
}
