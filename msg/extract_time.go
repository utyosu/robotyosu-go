package msg

import (
	"regexp"
	"strconv"
	"time"
)

var (
	parseTimeFuncs = []func(string) *timeInt{
		parseTimeHourMin,
		parseTimeHourMinJp,
		parseTimeHourMinHalfJp,
		parseTimeHourJp,
	}

	parseDateFuncs = []func(string) *dateInt{
		parseDateDayMonthYearSlash,
		parseDateDayMonthYearJp,
		parseDateDayMonthSlash,
		parseDateDayMonthJp,
		parseDateDayJp,
	}

	// varで定義すると実行開始時にコンパイルされるので高速に処理できる
	regexpDeleteWord                 = regexp.MustCompile(`[～-](\d+:\d+|\d+時)|(\d+:\d+|\d+時)まで`)
	regexpParseTimeHourMin           = regexp.MustCompile(`(\d+):(\d+)`)
	regexpParseTimeHourMinJp         = regexp.MustCompile(`(\d+)時(\d+)分`)
	regexpParseTimeHourMinHalfJp     = regexp.MustCompile(`(\d+)時半`)
	regexpParseTimeHourJp            = regexp.MustCompile(`(\d+)時`)
	regexpParseDateDayJp             = regexp.MustCompile(`(\d+)日`)
	regexpParseDateDayMonthJp        = regexp.MustCompile(`(\d+)月(\d+)日`)
	regexpParseDateDayMonthSlash     = regexp.MustCompile(`(\d+)/(\d+)`)
	regexpParseDateDayMonthYearJp    = regexp.MustCompile(`(\d+)年(\d+)月(\d+)日`)
	regexpParseDateDayMonthYearSlash = regexp.MustCompile(`(\d+)/(\d+)/(\d+)`)
)

type dateInt struct {
	year  int
	month int
	day   int
}

type timeInt struct {
	hour int
	min  int
}

func ParseTime(message string, now time.Time) *time.Time {
	// 終了時刻なら削除する
	message = regexpDeleteWord.ReplaceAllString(message, "")
	date := parseDateStringToInt(message)
	time := parseTimeStringToInt(message)
	if date != nil && time != nil {
		return date.toTimeWithTimeInt(time, now)
	} else if date == nil && time != nil {
		return time.toTime(now)
	} else if date != nil && time == nil {
		return date.toTimeWithTimeInt(&timeInt{hour: now.Hour(), min: now.Minute()}, now)
	}
	return nil
}

func parseTimeStringToInt(message string) *timeInt {
	for _, f := range parseTimeFuncs {
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
	return &timeInt{hour: hour, min: min}
}

// HH時MM分
func parseTimeHourMinJp(s string) *timeInt {
	c := regexpParseTimeHourMinJp.FindStringSubmatch(s)
	if len(c) < 3 {
		return nil
	}
	hour, _ := strconv.Atoi(c[1])
	min, _ := strconv.Atoi(c[2])
	return &timeInt{hour: hour, min: min}
}

// HH時半
func parseTimeHourMinHalfJp(s string) *timeInt {
	c := regexpParseTimeHourMinHalfJp.FindStringSubmatch(s)
	if len(c) < 2 {
		return nil
	}
	hour, _ := strconv.Atoi(c[1])
	return &timeInt{hour: hour, min: 30}
}

// HH時
func parseTimeHourJp(s string) *timeInt {
	c := regexpParseTimeHourJp.FindStringSubmatch(s)
	if len(c) < 2 {
		return nil
	}
	hour, _ := strconv.Atoi(c[1])
	return &timeInt{hour: hour, min: 0}
}

func parseDateStringToInt(message string) *dateInt {
	for _, f := range parseDateFuncs {
		if t := f(message); t != nil {
			return t
		}
	}
	return nil
}

// DD日
func parseDateDayJp(s string) *dateInt {
	c := regexpParseDateDayJp.FindStringSubmatch(s)
	if len(c) < 2 {
		return nil
	}
	day, _ := strconv.Atoi(c[1])
	return &dateInt{day: day}
}

// MM月DD日
func parseDateDayMonthJp(s string) *dateInt {
	c := regexpParseDateDayMonthJp.FindStringSubmatch(s)
	if len(c) < 3 {
		return nil
	}
	month, _ := strconv.Atoi(c[1])
	day, _ := strconv.Atoi(c[2])
	return &dateInt{month: month, day: day}
}

// MM/DD
func parseDateDayMonthSlash(s string) *dateInt {
	c := regexpParseDateDayMonthSlash.FindStringSubmatch(s)
	if len(c) < 3 {
		return nil
	}
	month, _ := strconv.Atoi(c[1])
	day, _ := strconv.Atoi(c[2])
	return &dateInt{month: month, day: day}
}

// YYYY年MM月DD日
func parseDateDayMonthYearJp(s string) *dateInt {
	c := regexpParseDateDayMonthYearJp.FindStringSubmatch(s)
	if len(c) < 4 {
		return nil
	}
	year, _ := strconv.Atoi(c[1])
	month, _ := strconv.Atoi(c[2])
	day, _ := strconv.Atoi(c[3])
	return &dateInt{year: year, month: month, day: day}
}

// YYYY/MM/DD
func parseDateDayMonthYearSlash(s string) *dateInt {
	c := regexpParseDateDayMonthYearSlash.FindStringSubmatch(s)
	if len(c) < 4 {
		return nil
	}
	year, _ := strconv.Atoi(c[1])
	month, _ := strconv.Atoi(c[2])
	day, _ := strconv.Atoi(c[3])
	return &dateInt{year: year, month: month, day: day}
}

func (t *timeInt) toTime(now time.Time) *time.Time {
	// 変な時間は変換しない
	if t.min < 0 || 60 <= t.min || t.hour < 0 || 30 <= t.hour {
		return nil
	}

	// 12時間表記で12時以降だと思われるものを24時間表記に変換する
	if t.hour <= 12 && t.hour < now.Hour() {
		t.hour += 12
		// それでも現在時刻より前なら翌日にする
		if t.hour < now.Hour() {
			t.hour += 12
		}
	}

	// 24時を超えていた場合は翌日扱いにする
	if 24 <= t.hour {
		result := time.Date(now.Year(), now.Month(), now.Day()+1, t.hour-24, t.min, 0, 0, now.Location())
		return &result
	}

	result := time.Date(now.Year(), now.Month(), now.Day(), t.hour, t.min, 0, 0, now.Location())
	return &result
}

func (d *dateInt) toTimeWithTimeInt(t *timeInt, now time.Time) *time.Time {
	// 変な時間は変換しない
	if t.min < 0 || 60 <= t.min || t.hour < 0 || 24 <= t.hour {
		return nil
	}

	if d.month == 0 {
		d.month = int(now.Month())
	}

	if d.year == 0 {
		d.year = now.Year()
	}

	result := time.Date(d.year, time.Month(d.month), d.day, t.hour, t.min, 0, 0, now.Location())
	return &result
}
