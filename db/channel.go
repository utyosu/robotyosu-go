package db

import (
	basic_errors "errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

const (
	defaultTimezone = "Asia/Tokyo"
)

type Channel struct {
	gorm.Model
	DiscordChannelId         int64
	DiscordGuildId           int64
	Recruitment              bool
	Timezone                 string
	Language                 string
	ReserveLimitTime         uint32
	ExpireDuration           uint32
	ExpireDurationForReserve uint32
	TwitterConfigId          uint32
}

func FindChannel(discordChannelId int64) (*Channel, error) {
	if r, found := mc.Get(getChannelCacheKey(discordChannelId)); found {
		return r.(*Channel), nil
	}

	channel := Channel{}
	if err := dbs.Take(&channel, "discord_channel_id=?", discordChannelId).Error; err != nil {
		if !basic_errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(err)
		}
	}

	mc.Set(
		getChannelCacheKey(channel.DiscordChannelId),
		&channel,
		cache.DefaultExpiration,
	)
	return &channel, nil
}

func FindOrCreateChannel(discordChannelId, discordGuildId int64) (*Channel, error) {
	channel, err := FindChannel(discordChannelId)
	if err != nil {
		return nil, err
	}
	if channel.ID == 0 {
		channel.DiscordChannelId = discordChannelId
		channel.DiscordGuildId = discordGuildId
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

func (c *Channel) IsEnabledRecruitment() bool {
	return c.ID != 0 && c.Recruitment
}

func (c *Channel) UpdateRecruitment(isRecruitment bool) error {
	c.Recruitment = isRecruitment
	err := dbs.Save(c).Error
	return errors.WithStack(err)
}

func (c *Channel) UpdateTimezone(timezone string) error {
	c.Timezone = timezone
	err := dbs.Save(c).Error
	return errors.WithStack(err)
}

func (c *Channel) UpdateLanguage(language string) error {
	c.Language = language
	err := dbs.Save(c).Error
	return errors.WithStack(err)
}

func (c *Channel) UpdateReserveLimitTime(reserveLimitTime uint32) error {
	c.ReserveLimitTime = reserveLimitTime
	err := dbs.Save(c).Error
	return errors.WithStack(err)
}

func (c *Channel) UpdateExpireDuration(expireDuration uint32) error {
	c.ExpireDuration = expireDuration
	err := dbs.Save(c).Error
	return errors.WithStack(err)
}

func (c *Channel) UpdateExpireDurationForReserve(expireDurationForReserve uint32) error {
	c.ExpireDurationForReserve = expireDurationForReserve
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

func getChannelCacheKey(discordChannelId int64) string {
	return fmt.Sprintf("channel/discord_channel_id=%v", discordChannelId)
}
