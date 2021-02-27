package db

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type TwitterConfig struct {
	gorm.Model
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
	Title             string
}

func FindTwitterConfig(id uint) (*TwitterConfig, error) {
	twitterConfig := TwitterConfig{}
	err := dbs.Take(&twitterConfig, "id=?", id).Error
	return &twitterConfig, errors.WithStack(err)
}
