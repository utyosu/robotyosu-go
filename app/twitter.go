package app

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/pkg/errors"
	"github.com/utyosu/robotyosu-go/db"
	"github.com/utyosu/robotyosu-go/i18n"
	"regexp"
	"strings"
)

type TwitterType int

const (
	TwitterTypeOpen TwitterType = iota
	TwitterTypeUpdate
	TwitterTypeClose
)

func newTwitterClient(c *db.TwitterConfig) *twitter.Client {
	config := oauth1.NewConfig(c.ConsumerKey, c.ConsumerSecret)
	token := oauth1.NewToken(c.AccessToken, c.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return twitter.NewClient(httpClient)
}

func tweet(c *db.Channel, r *db.Recruitment, t TwitterType) {
	if c.TwitterConfigId == 0 {
		// Twitter設定が存在しないときは何もしない
		return
	}

	twitterConfig, err := db.FindTwitterConfig(c.TwitterConfigId)
	if err != nil {
		// Tweetに失敗するだけならユーザーに通知しない
		postSlackWarning(err)
		return
	} else if twitterConfig.ID == 0 {
		postSlackWarning(
			fmt.Errorf(
				"TwitterConfig が見つかりません。\nChannelId: %v",
				c.ID,
			),
		)
		return
	}

	if r.TweetId == 0 && t != TwitterTypeOpen {
		postSlackWarning(
			fmt.Errorf(
				"募集開始ではないのにリプライ用のTweetIdが見つかりません。\nChannelId: %v\nRecruitmentId: %v",
				c.ID,
				r.ID,
			),
		)
		return
	}

	twitterClient := newTwitterClient(twitterConfig)
	var status *twitter.StatusUpdateParams
	if r.TweetId != 0 {
		status = &twitter.StatusUpdateParams{
			InReplyToStatusID: r.TweetId,
		}
	}

	msg := short140ForTwitter(buildTwitterMessage(twitterConfig, c, r, t))
	tweet, _, err := twitterClient.Statuses.Update(toTwitterSafe(msg), status)
	if err != nil {
		sendMessageT(c, "twitter_error")
		postSlackWarning(errors.WithStack(err))
		return
	}
	if err := r.UpdateTweetId(tweet.ID); err != nil {
		// DB更新は失敗してもツイートが成功しているので、ユーザーにエラーは出力しない
		postSlackWarning(errors.WithStack(err))
	}
	return
}

func buildTwitterMessage(twitterConfig *db.TwitterConfig, c *db.Channel, r *db.Recruitment, t TwitterType) string {
	switch t {
	case TwitterTypeOpen, TwitterTypeUpdate:
		memberNames := r.MemberNames()
		if len(memberNames) > 0 {
			return i18n.T(c.Language, "twitter_recruit", twitterConfig.Title, r.Title, r.AuthorName()) +
				"\n" +
				i18n.T(c.Language, "twitter_members", strings.Join(memberNames, ", "))
		} else {
			return fmt.Sprintf("%v\n%v by %v", twitterConfig.Title, r.Title, r.AuthorName())
		}
	case TwitterTypeClose:
		return i18n.T(c.Language, "twitter_close", twitterConfig.Title)
	}
	return ""
}

func toTwitterSafe(s string) string {
	return regexp.MustCompile(`@`).ReplaceAllString(s, "@ ")
}

func short140ForTwitter(s string) string {
	runes := []rune(s)
	if len(runes) <= 140 {
		return s
	}
	return string(runes[:137]) + "..."
}
