package config

import "os"

type DB struct {
	DB_host     string
	DB_name     string
	DB_user     string
	DB_password string
	DB_port     string
}

type AppConf struct {
	DB DB
}

func NewAppConf() AppConf {
	appConf := AppConf{
		DB: DB{
			DB_host:     os.Getenv("DB_host"),
			DB_name:     os.Getenv("DB_name"),
			DB_user:     os.Getenv("DB_user"),
			DB_password: os.Getenv("DB_password"),
			DB_port:     os.Getenv("DB_port"),
		},
	}

	return appConf
}
