package app

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/utyosu/robotyosu-go/db"
	"github.com/utyosu/robotyosu-go/i18n"
	"github.com/utyosu/robotyosu-go/msg"
	"golang.org/x/text/width"
)

var (
	regexpMention                 = regexp.MustCompile(`<@!?\d+>`)
	regexpOpenRecruitment         = regexp.MustCompile(`@(\d+)`)
	regexpFormatContentDeleteWord = regexp.MustCompile(`\r\n|\r|\n`)
)

func actionRecruitment(_ *discordgo.Session, m *discordgo.MessageCreate, channel *db.Channel, user *db.User) error {
	rawContent := regexpFormatContentDeleteWord.ReplaceAllString(m.Content, "")
	formattedContent := width.Fold.String(rawContent)
	switch {
	// 一覧
	case isContainKeywords(formattedContent, keywordsViewRecruitment):
		viewActiveRecruitments(channel)

	// 募集
	case regexpOpenRecruitment.MatchString(formattedContent):
		if haveMention(formattedContent) {
			return nil
		}
		timezone := channel.LoadLocation()
		now := time.Now().In(timezone)
		reserveAt := msg.ParseTime(formattedContent, now)
		capacity := uint32(getMatchRegexpNumber(formattedContent, regexpOpenRecruitment) + 1)
		recruitment, msg, err := db.InsertRecruitment(user, channel, rawContent, capacity, reserveAt)
		if err != nil {
			return err
		} else if recruitment == nil {
			sendMessage(m.ChannelID, msg)
			return nil
		}
		tweet(channel, recruitment, TwitterTypeOpen)
		if recruitment.ReserveAt != nil {
			sendMessageT(channel, "open_with_reserve", user.DisplayName(), recruitment.Label, recruitment.ReserveAtTime(channel.LoadLocation()))
		} else {
			sendMessageT(channel, "open", user.DisplayName(), recruitment.Label, recruitment.ExpireAtTime(channel.LoadLocation()))
		}
		viewActiveRecruitments(channel)

	// 参加
	case isContainKeywords(formattedContent, keywordsJoinRecruitment):
		recruitment, err := fetchRecruitmentWithMessage(formattedContent, channel)
		if err != nil {
			return err
		} else if recruitment == nil {
			return nil
		}
		if recruitment.IsParticipantsFull() {
			sendMessageT(channel, "not_join_because_full", recruitment.Label)
			viewActiveRecruitments(channel)
			return nil
		}
		if ok, err := recruitment.JoinParticipant(user, channel); err != nil {
			return err
		} else if ok {
			tweet(channel, recruitment, TwitterTypeUpdate)
			sendMessageT(channel, "join", user.DisplayName(), recruitment.Label)
			if recruitment.IsParticipantsFull() {
				if recruitment.IsPastReserveAt() {
					if err := recruitment.CloseRecruitment(); err != nil {
						return err
					}
					tweet(channel, recruitment, TwitterTypeClose)
					sendMessageT(channel, "gathered", recruitmentMentions(recruitment), recruitment.Label)
				} else {
					sendMessageT(channel, "gathered_reserved", recruitment.Label)
				}
			}
			viewActiveRecruitments(channel)
		}

	// 補欠参加
	case isContainKeywords(formattedContent, keywordsJoinRecruitmentAlternate):
		recruitment, err := fetchRecruitmentWithMessage(formattedContent, channel)
		if err != nil {
			return err
		} else if recruitment == nil {
			return nil
		}
		if ok, err := recruitment.JoinParticipantAlternate(user, channel); err != nil {
			return err
		} else if ok {
			sendMessageT(channel, "join_alternate", user.DisplayName(), recruitment.Label)
			viewActiveRecruitments(channel)
		}

	// キャンセル
	case isContainKeywords(formattedContent, keywordsCancelRecruitment):
		recruitment, err := fetchRecruitmentWithMessage(formattedContent, channel)
		if err != nil {
			return err
		} else if recruitment == nil {
			return nil
		}
		ok, err := recruitment.LeaveParticipant(user, channel)
		if err != nil {
			return err
		} else if ok {
			tweet(channel, recruitment, TwitterTypeUpdate)
			sendMessageT(channel, "leave", user.DisplayName(), recruitment.Label)
			viewActiveRecruitments(channel)
		}

	// 終了
	case isContainKeywords(formattedContent, keywordsCloseRecruitment):
		recruitment, err := fetchRecruitmentWithMessage(formattedContent, channel)
		if err != nil {
			return err
		} else if recruitment == nil {
			return nil
		}
		if err := recruitment.CloseRecruitment(); err != nil {
			return err
		}
		tweet(channel, recruitment, TwitterTypeClose)
		sendMessageT(channel, "closed", user.DisplayName(), recruitment.Label)
		viewActiveRecruitments(channel)

	// 復活
	case isContainKeywords(formattedContent, keywordsCloseResurrection):
		recruitment, err := db.ResurrectClosedRecruitment(channel)
		if err != nil {
			return err
		} else if recruitment != nil {
			sendMessageT(channel, "resurrection", recruitment.Label)
			viewActiveRecruitments(channel)
		}
	}
	return nil
}

func fetchRecruitmentWithMessage(content string, channel *db.Channel) (*db.Recruitment, error) {
	result, number := msg.ExtractNumber(trimMention(content))
	switch result {
	case msg.ExtractResultNotFoundNumber:
		return nil, nil
	case msg.ExtractResultMultipleNumber:
		sendMessageT(channel, "multiple_number")
		return nil, nil
	}
	recruitment, err := db.FetchActiveRecruitmentWithLabel(channel, number)
	if err != nil {
		return nil, err
	}
	if recruitment.ID == 0 {
		return nil, nil
	}
	return recruitment, nil
}

func haveMention(s string) bool {
	return regexpMention.MatchString(s)
}

func trimMention(s string) string {
	return regexpMention.ReplaceAllString(s, "")
}

func closeExpiredRecruitment() {
	channels, err := db.FetchAllChannels()
	if err != nil {
		slackWarning.Post(err)
		return
	}
	for _, channel := range channels {
		closed := false
		recruitments, err := db.FetchActiveRecruitments(channel)
		if err != nil {
			slackWarning.Post(err)
			return
		}
		for _, recruitment := range recruitments {
			if recruitment.IsPastExpireAt() {
				if err := recruitment.CloseRecruitment(); err != nil {
					slackWarning.Post(err)
					continue
				}
				tweet(channel, recruitment, TwitterTypeClose)
				sendMessageT(channel, "expired", recruitment.Label)
				closed = true
			}
		}
		if closed {
			viewActiveRecruitments(channel)
		}
	}
}

func notifyReservedRecruitmentOnTime() {
	now := time.Now()
	channels, err := db.FetchAllChannels()
	if err != nil {
		slackWarning.Post(err)
		return
	}
	for _, channel := range channels {
		recruitments, err := db.FetchActiveRecruitments(channel)
		if err != nil {
			slackWarning.Post(err)
			return
		}
		existNotified := false
		for _, recruitment := range recruitments {
			if notified, closed, err := recruitment.ProcessOnTime(now); err != nil {
				slackWarning.Post(err)
				return
			} else if notified {
				if closed {
					sendMessageT(channel, "close_reserved", recruitmentMentions(recruitment), recruitment.Label)
				} else {
					sendMessageT(channel, "reserve_on_time", recruitment.Label, recruitment.VacantSize())
				}
				existNotified = true
			}
		}
		if existNotified {
			viewActiveRecruitments(channel)
		}
	}
}

func viewActiveRecruitments(c *db.Channel) {
	recruitments, err := db.FetchActiveRecruitments(c)
	if err != nil {
		sendMessageT(c, "error")
		slackWarning.Post(err)
		return
	}
	m := "```\n"
	if len(recruitments) == 0 {
		m += i18n.T(c.Language, "no_recruitment")
	} else {
		for _, r := range recruitments {
			// 参加者が0人以下ならば表示しない
			if len(r.Participants) <= 0 {
				continue
			}

			// 参加メンバー表示
			memberNames := r.MemberNames(false)
			alternateMemberNames := r.MemberNames(true)

			m += i18n.T(c.Language, "recruit", r.Label, r.GetTitle(), r.AuthorName(), len(memberNames), r.Capacity-1) + "\n"

			if len(memberNames) > 0 {
				m += i18n.T(c.Language, "participants", strings.Join(memberNames, ", ")) + "\n"
			}
			if len(alternateMemberNames) > 0 {
				m += i18n.T(c.Language, "alternate_participants", strings.Join(alternateMemberNames, ", ")) + "\n"
			}
		}
	}
	m += "\n```"

	sendMessage(c.DiscordIdStr(), m)
}

func recruitmentMentions(r *db.Recruitment) string {
	var s = make([]string, len(r.Participants))
	for i, p := range r.Participants {
		s[i] = fmt.Sprintf("<@%v>", p.User.DiscordUserId)
	}
	return strings.Join(s, " ")
}
