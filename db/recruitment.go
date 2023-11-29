package db

import (
	basic_errors "errors"
	"time"

	"github.com/pkg/errors"
	"github.com/utyosu/robotyosu-go/db/encrypt"
	"github.com/utyosu/robotyosu-go/i18n"
	"gorm.io/gorm"
)

const (
	maxTitleRunes = 100

	// この時間以内なら終了済の募集であっても同じラベルを使わない
	ignoreLabelDuration = time.Hour * 6
)

type Recruitment struct {
	gorm.Model
	DiscordChannelId int64
	Label            uint32
	Title            string
	EncryptedTitle   []byte
	Capacity         uint32
	Active           bool
	Notified         bool
	TweetId          int64
	ReserveAt        *time.Time
	ExpireAt         *time.Time
	Participants     []Participant `gorm:"foreignkey:RecruitmentId"`
}

func (r *Recruitment) GetTitle() string {
	if len(r.EncryptedTitle) > 0 {
		return encrypt.DecryptString(r.EncryptedTitle)
	}
	return r.Title
}

// 指定ラベルの募集を取得する
// ユーザー名、ニックネームは取得しない
func FetchActiveRecruitmentWithLabel(channel *Channel, label int) (*Recruitment, error) {
	recruitment := &Recruitment{}
	err := dbs.
		Preload("Participants").
		Take(recruitment, "discord_channel_id=? AND active=? AND label=?", channel.DiscordChannelId, true, label).
		Error
	if basic_errors.Is(err, gorm.ErrRecordNotFound) {
		return recruitment, nil
	}
	return recruitment, errors.WithStack(err)
}

// 指定チャンネルの全募集を取得する
// ユーザー名、ニックネームも取得する
func FetchActiveRecruitments(channel *Channel) ([]*Recruitment, error) {
	recruitments := []*Recruitment{}
	err := dbs.
		Preload("Participants.User.Nickname", "discord_guild_id=?", channel.DiscordGuildId).
		Preload("Participants.User").
		Preload("Participants").
		Order("label ASC").
		Find(&recruitments, "discord_channel_id=? AND active=?", channel.DiscordChannelId, true).
		Error
	return recruitments, errors.WithStack(err)
}

func ResurrectClosedRecruitment(channel *Channel) (*Recruitment, error) {
	recruitment := &Recruitment{}
	err := dbs.
		Order("updated_at DESC").
		First(recruitment, "discord_channel_id=? AND active=? AND expire_at>?", channel.DiscordChannelId, false, time.Now()).
		Error
	if err != nil {
		return nil, errors.WithStack(err)
	} else if recruitment.ID != 0 {
		label, err := fetchEmptyLabel(channel)
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

func InsertRecruitment(user *User, channel *Channel, title string, capacity uint32, reserveAt *time.Time) (*Recruitment, string, error) {
	if len([]rune(title)) > maxTitleRunes {
		return nil, i18n.T(channel.Language, "too_long_title"), nil
	} else if capacity < 2 {
		return nil, i18n.T(channel.Language, "capacity_less"), nil
	} else if 4294967294 < capacity {
		return nil, i18n.T(channel.Language, "capacity_over"), nil
	}

	if reserveAt != nil {
		if reserveAt.Before(time.Now()) {
			reserveAt = nil
		} else {
			reserveAtDiff := uint32(reserveAt.Sub(time.Now()).Seconds())
			if channel.ReserveLimitTime != 0 && channel.ReserveLimitTime < reserveAtDiff {
				return nil, i18n.T(channel.Language, "reserve_limit_over"), nil
			}
		}
	}
	var expireAt time.Time
	if reserveAt != nil {
		expireAt = reserveAt.Add(time.Second * time.Duration(channel.ExpireDurationForReserve))
	} else {
		expireAt = time.Now().Add(time.Second * time.Duration(channel.ExpireDuration))
	}

	label, err := fetchEmptyLabel(channel)
	if err != nil {
		return nil, "", err
	}

	recruitment := &Recruitment{
		DiscordChannelId: channel.DiscordChannelId,
		Label:            label,
		EncryptedTitle:   encrypt.EncryptString(title),
		Capacity:         capacity,
		Active:           true,
		Notified:         false,
		ReserveAt:        reserveAt,
		ExpireAt:         &expireAt,
	}
	err = dbs.Create(recruitment).Error
	if err != nil {
		return nil, "", errors.WithStack(err)
	}
	if err := InsertParticipant(user, recruitment, false); err != nil {
		return nil, "", err
	}
	recruitment.Reload(channel)
	return recruitment, "", nil
}

func (r *Recruitment) CloseRecruitment() error {
	r.Active = false
	err := dbs.Save(r).Error
	return errors.WithStack(err)
}

func (r *Recruitment) JoinParticipant(user *User, channel *Channel) (bool, error) {
	for _, participant := range r.Participants {
		if participant.DiscordUserId == user.DiscordUserId {
			// 既に参加していれば失敗にする
			if !participant.Alternate {
				return false, nil
			}

			// 既に補欠参加していれば参加に切り替える
			if err := participant.Delete(); err != nil {
				return false, err
			}
			if err := InsertParticipant(user, r, false); err != nil {
				return false, err
			}
			return true, nil
		}
	}

	if err := InsertParticipant(user, r, false); err != nil {
		return false, err
	}
	r.Reload(channel)
	return true, nil
}

func (r *Recruitment) JoinParticipantAlternate(user *User, channel *Channel) (bool, error) {
	for i, participant := range r.Participants {
		if participant.DiscordUserId == user.DiscordUserId {

			// 最初の参加者(=主催者)のときは補欠参加できない
			if i == 0 {
				return false, nil
			}

			// 既に補欠参加していれば何もしない
			if participant.Alternate {
				return false, nil
			}

			// 既に参加していれば補欠に切り替える
			if err := participant.Delete(); err != nil {
				return false, err
			}
			if err := InsertParticipant(user, r, true); err != nil {
				return false, err
			}
			return true, nil
		}
	}

	// 補欠参加する
	if err := InsertParticipant(user, r, true); err != nil {
		return false, err
	}
	r.Reload(channel)
	return true, nil
}

func (r *Recruitment) LeaveParticipant(user *User, channel *Channel) (bool, error) {
	for _, p := range r.Participants {
		if p.RecruitmentId == uint32(r.ID) && p.DiscordUserId == user.DiscordUserId {
			if err := p.Delete(); err != nil {
				return false, err
			}
			r.Reload(channel)
			if r.ParticipantCount() == 0 {
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

func (r *Recruitment) ParticipantCount() int {
	var participantCount int
	for _, p := range r.Participants {
		if !p.Alternate {
			participantCount += 1
		}
	}
	return participantCount
}

func (r *Recruitment) IsParticipantsFull() bool {
	return int(r.Capacity) <= r.ParticipantCount()
}

func (r *Recruitment) VacantSize() int {
	return int(r.Capacity) - r.ParticipantCount()
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

func (r *Recruitment) Reload(channel *Channel) {
	// リロードは失敗しても影響が少ないのでエラーは無視する
	dbs.
		Preload("Participants.User.Nickname", "discord_guild_id=?", channel.DiscordGuildId).
		Preload("Participants.User").
		Preload("Participants").
		Take(r, "id=?", r.ID)
}

func (r *Recruitment) AuthorName() string {
	if len(r.Participants) <= 0 {
		return ""
	}
	return r.Participants[0].User.DisplayName()
}

func (r *Recruitment) MemberNames(alternate bool) []string {
	if len(r.Participants) <= 1 {
		return []string{}
	}
	memberSize := len(r.Participants) - 1
	names := make([]string, 0, memberSize)
	for _, p := range r.Participants[1:] {
		if p.Alternate == alternate {
			names = append(names, p.User.DisplayName())
		}
	}
	return names
}

func (r *Recruitment) ReserveAtTime(timezone *time.Location) string {
	if r.ReserveAt == nil {
		return ""
	}

	reserveAt := r.ReserveAt.In(timezone)
	now := time.Now().In(timezone)
	if reserveAt.Year() == now.Year() &&
		reserveAt.Month() == now.Month() &&
		reserveAt.Day() == now.Day() {
		return reserveAt.Format("15:04")
	}
	return reserveAt.Format("2006-01-02 15:04")
}

func (r *Recruitment) ExpireAtTime(timezone *time.Location) string {
	if r.ExpireAt == nil {
		return ""
	}
	return r.ExpireAt.In(timezone).Format("15:04")
}

func fetchEmptyLabel(channel *Channel) (uint32, error) {
	recruitments := []*Recruitment{}

	// 終了済の募集でもこの時間より後に作られた募集のラベルは使わない
	ignoreTime := time.Now().Add(-ignoreLabelDuration)

	err := dbs.
		Select("label").
		Find(&recruitments, "discord_channel_id = ? AND (active = ? OR created_at > ?)", channel.DiscordChannelId, true, ignoreTime).
		Error
	if err != nil {
		return 0, errors.WithStack(err)
	}

	var maxLabel uint32
	existLabelSet := map[uint32]struct{}{}
	for _, recruitment := range recruitments {
		existLabelSet[recruitment.Label] = struct{}{}
		if maxLabel < recruitment.Label {
			maxLabel = recruitment.Label
		}
	}

	for i := uint32(1); i < maxLabel; i++ {
		if _, ok := existLabelSet[i]; !ok {
			return i, nil
		}
	}

	return maxLabel + 1, nil
}
