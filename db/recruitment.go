package db

import (
	basic_errors "errors"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/utyosu/robotyosu-go/i18n"
	"time"
)

const (
	maxTitleLength = 100
)

type Recruitment struct {
	gorm.Model
	ChannelId    uint
	Label        uint
	Title        string
	Capacity     uint
	Active       bool
	Notified     bool
	TweetId      int64
	ReserveAt    *time.Time
	ExpireAt     *time.Time
	Participants []Participant `gorm:"foreignkey:RecruitmentId"`
}

func FetchActiveRecruitmentWithLabel(channelId uint, label int) (*Recruitment, error) {
	recruitment := &Recruitment{}
	err := dbs.Preload("Participants.User").Preload("Participants").Take(recruitment, "channel_id=? AND active=? AND label=?", channelId, true, label).Error
	if basic_errors.Is(err, gorm.ErrRecordNotFound) {
		return recruitment, nil
	}
	return recruitment, errors.WithStack(err)
}

func FetchActiveRecruitments(channelId uint) ([]*Recruitment, error) {
	recruitments := []*Recruitment{}
	err := dbs.Preload("Participants.User").Preload("Participants").Order("label ASC").Find(&recruitments, "channel_id=? AND active=?", channelId, true).Error
	return recruitments, errors.WithStack(err)
}

func ResurrectClosedRecruitment(channelId uint) (*Recruitment, error) {
	recruitment := &Recruitment{}
	err := dbs.Order("updated_at ASC").First(recruitment, "channel_id=? AND active=? AND expire_at>?", channelId, false, time.Now()).Error
	if err != nil {
		return nil, errors.WithStack(err)
	} else if recruitment.ID != 0 {
		label, err := fetchEmptyLabel(channelId)
		if err != nil {
			return nil, err
		}
		recruitment.Active = true
		recruitment.Label = label
		if err := dbs.Save(recruitment).Error; err != nil {
			return nil, errors.WithStack(err)
		}
		return recruitment, nil
	}
	return nil, nil
}

func InsertRecruitment(user *User, channel *Channel, title string, capacity uint, reserveAt *time.Time) (*Recruitment, string, error) {
	if len(title) > maxTitleLength {
		return nil, i18n.T(channel.Language, "too_long_title"), nil
	} else if capacity < 2 {
		return nil, i18n.T(channel.Language, "capacity_less"), nil
	}

	if reserveAt != nil && reserveAt.Before(time.Now()) {
		reserveAt = nil
	}
	var expireAt time.Time
	if reserveAt != nil {
		expireAt = reserveAt.Add(time.Minute * 30)
	} else {
		expireAt = time.Now().Add(time.Minute * 60)
	}

	label, err := fetchEmptyLabel(channel.ID)
	if err != nil {
		return nil, "", err
	}

	recruitment := &Recruitment{
		ChannelId: channel.ID,
		Label:     label,
		Title:     title,
		Capacity:  capacity,
		Active:    true,
		Notified:  false,
		ReserveAt: reserveAt,
		ExpireAt:  &expireAt,
	}
	err = dbs.Create(recruitment).Error
	if err != nil {
		return nil, "", errors.WithStack(err)
	}
	if err := InsertParticipant(user, recruitment); err != nil {
		return nil, "", err
	}
	recruitment.Reload()
	return recruitment, "", nil
}

func (r *Recruitment) CloseRecruitment() error {
	r.Active = false
	err := dbs.Save(r).Error
	return errors.WithStack(err)
}

func (r *Recruitment) JoinParticipant(user *User) (bool, error) {
	// 既に参加していれば失敗にする
	for _, participant := range r.Participants {
		if participant.UserId == user.ID {
			return false, nil
		}
	}

	if err := InsertParticipant(user, r); err != nil {
		return false, err
	}
	r.Reload()
	return true, nil
}

func (r *Recruitment) LeaveParticipant(user *User) (bool, error) {
	for _, p := range r.Participants {
		if p.RecruitmentId == r.ID && p.UserId == user.ID {
			if err := p.Delete(); err != nil {
				return false, err
			}
			r.Reload()
			if len(r.Participants) == 0 {
				if err := r.CloseRecruitment(); err != nil {
					return false, err
				}
			}
			return true, nil
		}
	}
	return false, nil
}

func (r *Recruitment) ProcessOnTime(now time.Time) (bool, bool, error) {
	// 予約なしは何もしない
	if r.ReserveAt == nil {
		return false, false, nil
	}

	var notified, closed bool

	// 未通知かつ時間を過ぎていれば通知する
	if !r.Notified && r.IsPastReserveAt() {
		if r.IsParticipantsFull() {
			// 集まっていればクローズする
			if err := r.CloseRecruitment(); err != nil {
				return false, false, err
			}
			closed = true
		}
		if err := r.NotifyRecruitment(); err != nil {
			return true, false, err
		}
		notified = true
	}
	return notified, closed, nil
}

// 予定時間が過ぎている、もしくは予定がなければtrue
// 予定時間が存在してまだ過ぎていないならfalse
func (r *Recruitment) IsPastReserveAt() bool {
	return r.ReserveAt == nil || r.ReserveAt.Before(time.Now())
}

func (r *Recruitment) IsPastExpireAt() bool {
	return r.ExpireAt == nil || r.ExpireAt.Before(time.Now())
}

func (r *Recruitment) IsParticipantsFull() bool {
	return int(r.Capacity) <= len(r.Participants)
}

func (r *Recruitment) VacantSize() int {
	return int(r.Capacity) - len(r.Participants)
}

func (r *Recruitment) NotifyRecruitment() error {
	r.Notified = true
	err := dbs.Save(r).Error
	return errors.WithStack(err)
}

func (r *Recruitment) UpdateTweetId(id int64) error {
	r.TweetId = id
	err := dbs.Save(r).Error
	return errors.WithStack(err)
}

func (r *Recruitment) Reload() {
	// リロードは失敗しても影響が少ないのでエラーは無視する
	dbs.Preload("Participants.User").Preload("Participants").Take(r, "id=?", r.ID)
}

func (r *Recruitment) AuthorName() string {
	if len(r.Participants) <= 0 {
		return ""
	}
	return r.Participants[0].User.Name
}

func (r *Recruitment) MemberNames() []string {
	if len(r.Participants) <= 1 {
		return []string{}
	}
	memberSize := len(r.Participants) - 1
	names := make([]string, memberSize, memberSize)
	for i, p := range r.Participants[1:] {
		names[i] = p.User.Name
	}
	return names
}

func (r *Recruitment) ReserveAtTime(timezone *time.Location) string {
	if r.ReserveAt == nil {
		return ""
	}
	return r.ReserveAt.In(timezone).Format("15:04")
}

func (r *Recruitment) ExpireAtTime(timezone *time.Location) string {
	if r.ExpireAt == nil {
		return ""
	}
	return r.ExpireAt.In(timezone).Format("15:04")
}

func fetchEmptyLabel(channelId uint) (uint, error) {
	recruitments, err := FetchActiveRecruitments(channelId)
	if err != nil {
		return 0, err
	}
	maxLabel := uint(1)
	for _, recruitment := range recruitments {
		if maxLabel <= recruitment.Label {
			maxLabel = recruitment.Label + 1
		}
	}
	return maxLabel, nil
}
