package db

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Participant struct {
	gorm.Model
	UserId        uint
	RecruitmentId uint
	User          User `gorm:"foreignkey:UserId"`
}

func InsertParticipant(user *User, recruitment *Recruitment) error {
	err := dbs.Create(&Participant{
		UserId:        user.ID,
		RecruitmentId: recruitment.ID,
	}).Error
	return errors.WithStack(err)
}

func (p *Participant) Delete() error {
	err := dbs.Delete(p).Error
	return errors.WithStack(err)
}
