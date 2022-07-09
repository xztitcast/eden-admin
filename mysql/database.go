package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

type Database struct {
	Orm *gorm.DB
}

func NewDatabase() *Database {
	dsn := "eden:y%^N&4by9k2XLD5S@tcp(101.33.211.193:3306)/g_sys?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Silent,
				Colorful:      true,
			}),
	})
	pool, err := db.DB()
	pool.SetMaxIdleConns(2)
	pool.SetMaxOpenConns(5)
	pool.SetConnMaxLifetime(time.Hour)
	if err != nil {
		panic(err)
	}
	return &Database{
		Orm: db,
	}
}
