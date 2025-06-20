package migrate

import (
	"fmt"
	"log"
	"petstore/config"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

func Migration(cfg config.AppConf) {
	str := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.DB.DB_user, cfg.DB.DB_password, cfg.DB.DB_host, cfg.DB.DB_port, cfg.DB.DB_name)
	m, err := migrate.New(
		"file://.",
		str,
	)
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Ошибка применения миграции: %v", err)
	}

	log.Println("Миграции успешно применены")
}
