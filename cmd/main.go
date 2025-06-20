package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"petstore/config"
	"petstore/internal/db"
	"petstore/internal/logs"
	"petstore/internal/migrate"
	"petstore/internal/route"
	"petstore/run"

	_ "petstore/internal/docs"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// @title Petstore
// @description Зоомагазин, который работает с тремя основными сущностями: пользователь, животное и магазин.
// @description Здесь доступны CRUD-операции над всеми сущностями, мягкая работа с пользователем (проставление дат действий), авторизация и выход из системы с помощью чёрного списка, загрузка изображения и сохранение локально на сервере. Это самые интересные возможности, попробуйте сами
// @version 1.0

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Ошибка чтения .env файла")
	}

	// переменные окружения
	cfg := config.NewAppConf()
	// подключение к бд
	db_conn := db.NewConnect(cfg)
	// настройка логгера
	logger := logs.NewLogger()
	// настройка миграций
	migrate.Migration(cfg)

	// user слои
	userRespond, middleAuth := run.NewModulesUser(db_conn, logger)
	// pet слои
	petRespond := run.NewModulesPet(db_conn, logger)
	// order слои
	orderRespond := run.NewModulesOrder(db_conn, logger)

	server := http.Server{
		Handler:      route.HandlerPetStore(middleAuth, userRespond, petRespond, orderRespond),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	chClose := make(chan os.Signal, 1)
	signal.Notify(chClose, syscall.SIGINT, syscall.SIGTERM)

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Fatal("Ошибка транспортного уровня", zap.String("err", err.Error()))
	}
	defer listen.Close()

	go func() {
		err := server.Serve(listen)
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("Ошибка запуска сервера", zap.String("err", err.Error()))
		}
	}()
	logger.Info("Сервер запущен", zap.String("addr", listen.Addr().String()))

	<-chClose
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		logger.Warn("Сервер некорректно завершил работу", zap.String("err", err.Error()))
		return
	}
	close(chClose)

	logger.Info("Сервер остановлен Graceful Shutdown")
}
