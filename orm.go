package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sync"
)

var DB *gorm.DB
var once sync.Once

func InitOrm() *gorm.DB {

	once.Do(func() {
		var err error
		DB, err = gorm.Open(sqlite.Open("./wallet.db"), &gorm.Config{
			SkipDefaultTransaction: true,
		})
		if err != nil {
			panic("failed to connect database")
		}

		if err := DB.AutoMigrate(&Gift{}, &User{}, &LogGift{}); err != nil {
			panic("failed to migrate database")
		}
		firstGift := DB.First(&Gift{}, "code = ?", "xxx")
		if firstGift.RowsAffected == 0 {
			DB.Create(&Gift{
				Code:  "xxx",
				Count: 1000,
				Price: 1000000,
			})
		}

	})
	return DB
}
