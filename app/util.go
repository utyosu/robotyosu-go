package app

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/utyosu/robotyosu-go/slack"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func isContainKeywords(m string, keywords []string) bool {
	for _, k := range keywords {
		if strings.Contains(m, k) {
			return true
		}
	}
	return false
}

func getMatchRegexpNumber(s string, r *regexp.Regexp) int {
	c := r.FindStringSubmatch(s)
	if len(c) < 2 {
		return 0
	}
	i, _ := strconv.Atoi(c[1])
	return i
}

func getMatchRegexpString(s string, r *regexp.Regexp) string {
	c := r.FindStringSubmatch(s)
	if len(c) < 2 {
		return ""
	}
	return c[1]
}

func doFuncSchedule(f func(), interval time.Duration) *time.Ticker {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			f()
		}
	}()
	return ticker
}

func NotifySlackWhenPanic(info string) {
	if err := recover(); err != nil {
		slack.PostSlackAlert(fmt.Sprintf("panic: %v\ninfo: %v", err, info))
	}
}

func messageInformation(s *discordgo.Session, m *discordgo.MessageCreate) string {
	return fmt.Sprintf(
		"session: %v\nmessage: %v",
		s,
		m,
	)
}
