package msg

import (
	"regexp"
	"strconv"
	"time"
)

var (
	parseFuncs = []func(string) *timeInt{
		parseTimeHourMin,
		parseTimeHourMinJp,
		parseTimeHourMinHalfJp,
		parseTimeHourJp,
	}

	// varで定義すると実行開始時にコンパイルされるので高速に処理できる
	regexpDeleteWord             = regexp.MustCompile(`[～-](\d+:\d+|\d+時)|(\d+:\d+|\d+時)まで`)
	regexpParseTimeHourMin       = regexp.MustCompile(`(\d+):(\d+)`)
	regexpParseTimeHourMinJp     = regexp.MustCompile(`(\d+)時(\d+)分`)
	regexpParseTimeHourMinHalfJp = regexp.MustCompile(`(\d+)時半`)
	regexpParseTimeHourJp        = regexp.MustCompile(`(\d+)時`)
)

type timeInt struct {
	hour int
	min  int
}

func ParseTime(message string, now time.Time) *time.Time {
	// 終了時刻なら削除する
	message = regexpDeleteWord.ReplaceAllString(message, "")
	if t := parseTimeStringToInt(message); t != nil {
		return t.toTime(now)
	}
	return nil
}

func parseTimeStringToInt(message string) *timeInt {
	for _, f := range parseFuncs {
		if t := f(message); t != nil {
			return t
		}
	}
	return nil
}

// HH:MM
func parseTimeHourMin(s string) *timeInt {
	c := regexpParseTimeHourMin.FindStringSubmatch(s)
	if len(c) < 3 {
		return nil
	}
	hour, _ := strconv.Atoi(c[1])
	min, _ := strconv.Atoi(c[2])
	return &timeInt{hour, min}
}

// HH時MM分
func parseTimeHourMinJp(s string) *timeInt {
	c := regexpParseTimeHourMinJp.FindStringSubmatch(s)
	if len(c) < 3 {
		return nil
	}
	hour, _ := strconv.Atoi(c[1])
	min, _ := strconv.Atoi(c[2])
	return &timeInt{hour, min}
}

// HH時半
func parseTimeHourMinHalfJp(s string) *timeInt {
	c := regexpParseTimeHourMinHalfJp.FindStringSubmatch(s)
	if len(c) < 2 {
		return nil
	}
	hour, _ := strconv.Atoi(c[1])
	return &timeInt{hour, 30}
}

// HH時
func parseTimeHourJp(s string) *timeInt {
	c := regexpParseTimeHourJp.FindStringSubmatch(s)
	if len(c) < 2 {
		return nil
	}
	hour, _ := strconv.Atoi(c[1])
	return &timeInt{hour, 0}
}

func (t *timeInt) toTime(now time.Time) *time.Time {
	// 変な時間は変換しない
	if t.min < 0 || 60 <= t.min || t.hour < 0 || 30 <= t.hour {
		return nil
	}

	// 12時間表記で12時以降だと思われるものを24時間表記に変換する
	if t.hour <= 12 && t.hour < now.Hour() {
		t.hour += 12
	}

	// 24時を超えていた場合は翌日扱いにする
	if 24 <= t.hour {
		result := time.Date(now.Year(), now.Month(), now.Day()+1, t.hour-24, t.min, 0, 0, now.Location())
		return &result
	}

	result := time.Date(now.Year(), now.Month(), now.Day(), t.hour, t.min, 0, 0, now.Location())
	return &result
}
