package main

import (
	"github.com/utyosu/robotyosu-go/app"
	"github.com/utyosu/robotyosu-go/db"
	"time"
)

func init() {
	if loc, err := time.LoadLocation("Asia/Tokyo"); err == nil {
		time.Local = loc
	}
}

func main() {
	defer app.NotifySlackWhenPanic("main")
	db.ConnectDb()
	app.Start()
	return
}
