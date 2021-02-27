package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/utyosu/robotyosu-go/env"
)

var (
	dbs *gorm.DB
)

func ConnectDb() *gorm.DB {
	var err error
	dbs, err = gorm.Open(
		env.DbDriver,
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true",
			env.DbUser,
			env.DbPassword,
			env.DbHost,
			env.DbPort,
			env.DbName,
		),
	)

	if err != nil {
		panic(err)
	}

	return dbs
}
