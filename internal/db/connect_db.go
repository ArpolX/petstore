package db

import (
	"fmt"
	"log"
	"os"
	"petstore/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewConnect(cfg config.AppConf) *gorm.DB {
	strConnect := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.DB.DB_host, cfg.DB.DB_port, cfg.DB.DB_user, cfg.DB.DB_name, cfg.DB.DB_password)

	db, err := gorm.Open(postgres.Open(strConnect), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\n", log.Flags()),
			logger.Config{
				LogLevel: logger.Info,
				Colorful: true,
			},
		),
	})
	if err != nil {
		log.Fatalf("Ошибка подключения к бд: %v", err)
	}

	return db
}
