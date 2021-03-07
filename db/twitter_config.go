package db

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TwitterConfig struct {
	gorm.Model
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
	Title             string
}

func FindTwitterConfig(id uint32) (*TwitterConfig, error) {
	twitterConfig := TwitterConfig{}
	err := dbs.Take(&twitterConfig, "id=?", id).Error
	return &twitterConfig, errors.WithStack(err)
}
