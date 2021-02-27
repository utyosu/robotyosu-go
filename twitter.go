package main

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/utyosu/robotyosu-go/db"
	"github.com/utyosu/robotyosu-go/i18n"
	"log"
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
		log.Println("[Error] TwitterConfig が見つかりません。")
		return
	}

	if r.TweetId == 0 && t != TwitterTypeOpen {
		log.Println("[Error] 募集開始ではないのにリプライ用のTweetIdが見つかりません。")
		return
	}

	twitterClient := newTwitterClient(twitterConfig)
	var status *twitter.StatusUpdateParams
	if r.TweetId != 0 {
		status = &twitter.StatusUpdateParams{
			InReplyToStatusID: r.TweetId,
		}
	}

	msg := buildTwitterMessage(twitterConfig, c, r, t)
	tweet, _, err := twitterClient.Statuses.Update(toTwitterSafe(msg), status)
	if err != nil {
		log.Println(err)
		return
	}
	if err := r.UpdateTweetId(tweet.ID); err != nil {
		// DB更新は失敗してもツイートが成功しているので、ユーザーにエラーは出力しない
		postSlackWarning(err)
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
