package app

import (
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

func NotifySlackWhenPanic(p ...interface{}) {
	if err := recover(); err != nil {
		slack.PostSlackAlert(p...)
	}
}
