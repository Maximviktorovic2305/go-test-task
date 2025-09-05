package database

import (
	"log"

	"effective-mobile-subscription/internal/models"

	"gorm.io/gorm"
)

// Migrate выполняет миграции базы данных
func Migrate(db *gorm.DB) {
	log.Println("Выполнение миграций...")

	// Мигрировать модель Subscription
	err := db.AutoMigrate(&models.Subscription{})
	if err != nil {
		log.Fatal("Не удалось выполнить миграцию базы данных:", err)
	}

	log.Println("Миграции успешно завершены")
}
