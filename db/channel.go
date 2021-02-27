package db

import (
	basic_errors "errors"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"log"
	"strconv"
	"time"
)

const (
	defaultTimezone = "Asia/Tokyo"
)

type Channel struct {
	gorm.Model
	DiscordChannelId int64
	Recruitment      bool
	Timezone         string
	Language         string
	TwitterConfigId  uint
}

func FindChannel(discordChannelId int64) (*Channel, error) {
	channel := Channel{}
	if err := dbs.Take(&channel, "discord_channel_id=?", discordChannelId).Error; err != nil {
		if basic_errors.Is(err, gorm.ErrRecordNotFound) {
			return &channel, nil
		}
		return nil, errors.WithStack(err)
	}
	return &channel, nil
}

func FindOrCreateChannel(discordChannelId int64) (*Channel, error) {
	channel, err := FindChannel(discordChannelId)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if channel.ID == 0 {
		channel.DiscordChannelId = discordChannelId
		channel.Timezone = defaultTimezone
		if err := dbs.Create(&channel).Error; err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return channel, nil
}

func FetchAllChannels() ([]*Channel, error) {
	channels := []*Channel{}
	err := dbs.Find(&channels, "recruitment=?", true).Error
	return channels, errors.WithStack(err)
}

func (c *Channel) IsValidAnyFunction() bool {
	return c.ID != 0 && c.Recruitment
}

func (c *Channel) UpdateChannelRecruitment(isRecruitment bool) error {
	c.Recruitment = isRecruitment
	err := dbs.Save(c).Error
	return errors.WithStack(err)
}

func (c *Channel) UpdateChannelTimezone(timezone string) error {
	c.Timezone = timezone
	err := dbs.Save(c).Error
	return errors.WithStack(err)
}

func (c *Channel) UpdateChannelLanguage(language string) error {
	c.Language = language
	err := dbs.Save(c).Error
	return errors.WithStack(err)
}

func (c *Channel) DiscordIdStr() string {
	return strconv.FormatInt(c.DiscordChannelId, 10)
}

func (c *Channel) LoadLocation() *time.Location {
	timezone, err := time.LoadLocation(c.Timezone)
	if err != nil {
		log.Println(err)
		timezone, _ = time.LoadLocation(defaultTimezone)
	}
	return timezone
}
