package services

import (
	"effective-mobile-subscription/internal/models"
	"effective-mobile-subscription/pkg/utils"

	"gorm.io/gorm"
)

// SubscriptionService обрабатывает бизнес-логику для подписок
type SubscriptionService struct {
	db *gorm.DB
}

// NewSubscriptionService создает новый сервис подписок
func NewSubscriptionService(db *gorm.DB) *SubscriptionService {
	return &SubscriptionService{db: db}
}

// CreateSubscription создает новую подписку
func (s *SubscriptionService) CreateSubscription(subscription *models.Subscription) error {
	// Убедиться, что StartDate установлена на первый день месяца
	subscription.StartDate = utils.GetFirstDayOfMonth(subscription.StartDate)

	// Если EndDate предоставлена, установить ее на последний день месяца
	if subscription.EndDate != nil {
		endDate := utils.GetLastDayOfMonth(*subscription.EndDate)
		subscription.EndDate = &endDate
	}

	return s.db.Create(subscription).Error
}

// GetSubscription получает подписку по ID
func (s *SubscriptionService) GetSubscription(id uint) (*models.Subscription, error) {
	var subscription models.Subscription
	err := s.db.First(&subscription, id).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// UpdateSubscription обновляет существующую подписку
func (s *SubscriptionService) UpdateSubscription(id uint, subscription *models.Subscription) error {
	// Убедиться, что StartDate установлена на первый день месяца
	subscription.StartDate = utils.GetFirstDayOfMonth(subscription.StartDate)

	// Если EndDate предоставлена, установить ее на последний день месяца
	if subscription.EndDate != nil {
		endDate := utils.GetLastDayOfMonth(*subscription.EndDate)
		subscription.EndDate = &endDate
	}

	// Обновить только предоставленные поля
	return s.db.Model(&models.Subscription{}).Where("id = ?", id).Updates(subscription).Error
}

// DeleteSubscription удаляет подписку по ID
func (s *SubscriptionService) DeleteSubscription(id uint) error {
	return s.db.Delete(&models.Subscription{}, id).Error
}

// ListSubscriptions получает список подписок с опциональными фильтрами и пагинацией
func (s *SubscriptionService) ListSubscriptions(page, limit int, userID, serviceName string) ([]models.Subscription, int64, error) {
	var subscriptions []models.Subscription
	var total int64

	// Вычислить смещение для пагинации
	offset := (page - 1) * limit

	// Построить запрос с фильтрами
	query := s.db.Model(&models.Subscription{})

	// Применить фильтры
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	// Получить общее количество
	query.Count(&total)

	// Получить результаты с пагинацией
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&subscriptions).Error
	if err != nil {
		return nil, 0, err
	}

	return subscriptions, total, nil
}

// CalculateTotalCost вычисляет общую стоимость подписок с опциональными фильтрами
func (s *SubscriptionService) CalculateTotalCost(userID, serviceName, from, to string) (int, error) {
	var totalCost int64

	// Построить запрос
	query := s.db.Model(&models.Subscription{}).Select("SUM(price)")

	// Применить фильтры
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	// Применить фильтры по диапазону дат
	if from != "" {
		fromDate, err := utils.ParseMonthYear(from)
		if err != nil {
			return 0, err
		}
		// Установить на первый день месяца
		fromDate = utils.GetFirstDayOfMonth(fromDate)
		query = query.Where("start_date >= ?", fromDate)
	}

	if to != "" {
		toDate, err := utils.ParseMonthYear(to)
		if err != nil {
			return 0, err
		}
		// Установить на последний день месяца
		toDate = utils.GetLastDayOfMonth(toDate)
		query = query.Where("start_date <= ?", toDate)
	}

	// Выполнить запрос
	err := query.Scan(&totalCost).Error
	if err != nil {
		return 0, err
	}

	return int(totalCost), nil
}

// applyFilters применяет области GORM для фильтрации подписок
func (s *SubscriptionService) applyFilters(query *gorm.DB, userID, serviceName string) *gorm.DB {
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}
	return query
}
