package services

import (
	"effective-mobile-subscription/internal/models"
	"effective-mobile-subscription/internal/repository"
	"effective-mobile-subscription/pkg/utils"
)

// SubscriptionService обрабатывает бизнес-логику для подписок
type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

// NewSubscriptionService создает новый сервис подписок
func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

// CreateSubscription создает новую подписку
func (s *SubscriptionService) CreateSubscription(subscription *models.Subscription) error {
	// Убедимся, что StartDate установлена на первый день месяца
	subscription.StartDate = utils.GetFirstDayOfMonth(subscription.StartDate)

	// Если EndDate предоставлена, установить ее на последний день месяца
	if subscription.EndDate != nil {
		endDate := utils.GetLastDayOfMonth(*subscription.EndDate)
		subscription.EndDate = &endDate
	}

	return s.repo.Create(subscription)
}

// GetSubscription получает подписку по ID
func (s *SubscriptionService) GetSubscription(id uint) (*models.Subscription, error) {
	return s.repo.GetByID(id)
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

	return s.repo.Update(id, subscription)
}

// DeleteSubscription удаляет подписку по ID
func (s *SubscriptionService) DeleteSubscription(id uint) error {
	return s.repo.Delete(id)
}

// ListSubscriptions получает список подписок с опциональными фильтрами и пагинацией
func (s *SubscriptionService) ListSubscriptions(page, limit int, userID, serviceName string) ([]models.Subscription, int64, error) {
	// Вычислить смещение для пагинации
	offset := (page - 1) * limit

	// Получить результаты с пагинацией
	subscriptions, err := s.repo.List(offset, limit, userID, serviceName)
	if err != nil {
		return nil, 0, err
	}

	// Получить общее количество
	total, err := s.repo.Count(userID, serviceName)
	if err != nil {
		return nil, 0, err
	}

	return subscriptions, total, nil
}

// CalculateTotalCost вычисляет общую стоимость подписок с опциональными фильтрами
func (s *SubscriptionService) CalculateTotalCost(userID, serviceName, from, to string) (int, error) {
	// Применить фильтры по диапазону дат
	var startDate, endDate any
	if from != "" {
		fromDate, err := utils.ParseMonthYear(from)
		if err != nil {
			return 0, err
		}
		// Установить на первый день месяца
		fromDate = utils.GetFirstDayOfMonth(fromDate)
		startDate = fromDate
	}

	if to != "" {
		toDate, err := utils.ParseMonthYear(to)
		if err != nil {
			return 0, err
		}
		// Установить на последний день месяца
		toDate = utils.GetLastDayOfMonth(toDate)
		endDate = toDate
	}

	// Выполнить запрос
	totalCost, err := s.repo.CalculateTotalCost(userID, serviceName, startDate, endDate)
	if err != nil {
		return 0, err
	}

	return int(totalCost), nil
}
