package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"wallet-api/config"
	"wallet-api/internal/handler"
	"wallet-api/internal/middleware"
	"wallet-api/internal/repository"
	"wallet-api/internal/service"
	"wallet-api/utils/logger"

	_ "github.com/lib/pq"
)

func main() {
	logger.Init()

	if err := config.NewConf(); err != nil {
		logger.GlobalLogger.Error("Ошибка загрузки конфигурации: %v", err)
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Cnf.PgHost,
		config.Cnf.PgPort,
		config.Cnf.PgUser,
		config.Cnf.PgPassword,
		config.Cnf.PgDbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.GlobalLogger.Error("Ошибка подключения к БД: %v", err)
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(config.Cnf.MaxConnections)
	db.SetMaxIdleConns(config.Cnf.MaxConnections / 2)
	db.SetConnMaxLifetime(time.Minute * 5)

	if err = db.Ping(); err != nil {
		logger.GlobalLogger.Error("Ошибка проверки соединения с БД: %v", err)
		log.Fatal(err)
	}

	walletRepo := repository.NewWalletRepository(db)
	walletService := service.NewWalletService(walletRepo)
	walletHandler := handler.NewWalletHandler(walletService)

	http.HandleFunc("/api/v1/wallet", middleware.LoggingMiddleware(walletHandler.HandleWalletOperation))
	http.HandleFunc("/api/v1/wallets/", middleware.LoggingMiddleware(walletHandler.HandleGetWallet))

	logger.GlobalLogger.Info("Сервер запущен на порту %s", config.Cnf.HttpPort)
	log.Fatal(http.ListenAndServe(":"+config.Cnf.HttpPort, nil))
}
