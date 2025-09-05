package repository

import (
	"effective-mobile-subscription/internal/models"

	"gorm.io/gorm"
)

// SubscriptionRepository обрабатывает операции с базой данных для подписок
type SubscriptionRepository struct {
	db *gorm.DB
}

// NewSubscriptionRepository создает новый репозиторий подписок
func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create создает новую подписку в базе данных
func (r *SubscriptionRepository) Create(subscription *models.Subscription) error {
	return r.db.Create(subscription).Error
}

// GetByID получает подписку по её ID
func (r *SubscriptionRepository) GetByID(id uint) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.First(&subscription, id).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// Update обновляет существующую подписку
func (r *SubscriptionRepository) Update(id uint, subscription *models.Subscription) error {
	// Обновить только предоставленные поля
	return r.db.Model(&models.Subscription{}).Where("id = ?", id).Updates(subscription).Error
}

// Delete удаляет подписку по её ID
func (r *SubscriptionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Subscription{}, id).Error
}

// List получает подписки с опциональными фильтрами и пагинацией
func (r *SubscriptionRepository) List(offset, limit int, userID, serviceName string) ([]models.Subscription, error) {
	var subscriptions []models.Subscription

	// Построить запрос с фильтрами
	query := r.db.Model(&models.Subscription{})

	// Применить фильтры
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	// Получить результаты с пагинацией
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&subscriptions).Error
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// Count возвращает общее количество подписок с опциональными фильтрами
func (r *SubscriptionRepository) Count(userID, serviceName string) (int64, error) {
	var total int64

	// Построить запрос с фильтрами
	query := r.db.Model(&models.Subscription{})

	// Применить фильтры
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	// Получить количество
	err := query.Count(&total).Error
	if err != nil {
		return 0, err
	}

	return total, nil
}

// CalculateTotalCost вычисляет общую стоимость подписок с опциональными фильтрами
func (r *SubscriptionRepository) CalculateTotalCost(userID, serviceName string, startDate, endDate any) (int64, error) {
	var totalCost int64

	// Построить запрос
	query := r.db.Model(&models.Subscription{}).Select("SUM(price)")

	// Применить фильтры
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	// Применить фильтры по диапазону дат, если они предоставлены
	if startDate != nil {
		query = query.Where("start_date >= ?", startDate)
	}

	if endDate != nil {
		query = query.Where("start_date <= ?", endDate)
	}

	// Выполнить запрос
	err := query.Scan(&totalCost).Error
	if err != nil {
		return 0, err
	}

	return totalCost, nil
}
