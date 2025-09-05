package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"effective-mobile-subscription/config"
	"effective-mobile-subscription/internal/database"
	"effective-mobile-subscription/internal/routes"
	"effective-mobile-subscription/pkg/middleware"
	"effective-mobile-subscription/pkg/utils"
)

func main() {
	// Загрузить конфигурацию
	cfg := config.LoadConfig()

	// Создать логгер
	logger := utils.NewLogger()
	logger.Info("Запуск сервиса подписок")

	// Подключиться к базе данных
	db := database.ConnectDB(cfg)

	// Выполнить миграции
	database.Migrate(db)

	// Настроить маршруты
	router := routes.SetupRoutes(db)

	// Добавить middleware
	router.Use(middleware.ErrorMiddleware(logger))

	// Создать HTTP сервер
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Запустить сервер в горутине
	go func() {
		logger.Info("Сервер запускается", "адрес", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Не удалось запустить сервер: %v", err)
		}
	}()

	// Дождаться сигнала прерывания для корректного завершения работы сервера
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Завершение работы сервера...")

	// Контекст используется для информирования сервера о том, что у него есть 5 секунд для завершения
	// текущий обрабатываемый запрос
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Сервер принудительно завершен", "ошибка", err)
	}

	logger.Info("Сервер завершен")
}
