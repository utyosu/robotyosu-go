package main

import (
	"github.com/utyosu/robotyosu-go/app"
	"time"
)

func init() {
	if loc, err := time.LoadLocation("Asia/Tokyo"); err == nil {
		time.Local = loc
	}
}

func main() {
	defer app.NotifySlackWhenPanic("main")
	app.Start()
	return
}
