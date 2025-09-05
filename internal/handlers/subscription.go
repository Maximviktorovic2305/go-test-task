package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"effective-mobile-subscription/internal/models"
	"effective-mobile-subscription/internal/services"
	"effective-mobile-subscription/pkg/utils"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// SubscriptionHandler обрабатывает HTTP запросы для подписок
type SubscriptionHandler struct {
	service *services.SubscriptionService
}

// NewSubscriptionHandler создает новый обработчик подписок
func NewSubscriptionHandler(service *services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// CreateSubscription создает новую подписку
func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServiceName string `json:"service_name"`
		Price       int    `json:"price"`
		UserID      string `json:"user_id"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date,omitempty"`
	}

	// Декодировать тело запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверное тело запроса", http.StatusBadRequest)
		return
	}

	// Разобрать дату начала
	startDate, err := utils.ParseMonthYear(req.StartDate)
	if err != nil {
		http.Error(w, "Неверный формат даты начала, ожидается ММ-ГГГГ", http.StatusBadRequest)
		return
	}

	// Разобрать дату окончания, если предоставлена
	var endDate *time.Time
	if req.EndDate != "" {
		end, err := utils.ParseMonthYear(req.EndDate)
		if err != nil {
			http.Error(w, "Неверный формат даты окончания, ожидается ММ-ГГГГ", http.StatusBadRequest)
			return
		}
		endDate = &end
	}

	// Создать модель подписки
	subscription := &models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	// Создать подписку в базе данных
	if err := h.service.CreateSubscription(subscription); err != nil {
		http.Error(w, "Не удалось создать подписку", http.StatusInternalServerError)
		return
	}

	// Вернуть созданную подписку
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subscription)
}

// GetSubscription получает подписку по ID
func (h *SubscriptionHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	// Получить ID из параметров URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID подписки", http.StatusBadRequest)
		return
	}

	// Получить подписку из базы данных
	subscription, err := h.service.GetSubscription(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Подписка не найдена", http.StatusNotFound)
			return
		}
		http.Error(w, "Не удалось получить подписку", http.StatusInternalServerError)
		return
	}

	// Вернуть подписку
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)
}

// UpdateSubscription обновляет существующую подписку
func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	// Получить ID из параметров URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID подписки", http.StatusBadRequest)
		return
	}

	var req struct {
		ServiceName string `json:"service_name,omitempty"`
		Price       int    `json:"price,omitempty"`
		UserID      string `json:"user_id,omitempty"`
		StartDate   string `json:"start_date,omitempty"`
		EndDate     string `json:"end_date,omitempty"`
	}

	// Декодировать тело запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверное тело запроса", http.StatusBadRequest)
		return
	}

	// Разобрать дату начала, если предоставлена
	var startDate *time.Time
	if req.StartDate != "" {
		start, err := utils.ParseMonthYear(req.StartDate)
		if err != nil {
			http.Error(w, "Неверный формат даты начала, ожидается ММ-ГГГГ", http.StatusBadRequest)
			return
		}
		startDate = &start
	}

	// Разобрать дату окончания, если предоставлена
	var endDate *time.Time
	if req.EndDate != "" {
		end, err := utils.ParseMonthYear(req.EndDate)
		if err != nil {
			http.Error(w, "Неверный формат даты окончания, ожидается ММ-ГГГГ", http.StatusBadRequest)
			return
		}
		endDate = &end
	}

	// Создать модель подписки с обновленными полями
	subscription := &models.Subscription{}
	if req.ServiceName != "" {
		subscription.ServiceName = req.ServiceName
	}
	if req.Price != 0 {
		subscription.Price = req.Price
	}
	if req.UserID != "" {
		subscription.UserID = req.UserID
	}
	if startDate != nil {
		subscription.StartDate = *startDate
	}
	if endDate != nil {
		subscription.EndDate = endDate
	}

	// Обновить подписку в базе данных
	if err := h.service.UpdateSubscription(uint(id), subscription); err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Подписка не найдена", http.StatusNotFound)
			return
		}
		http.Error(w, "Не удалось обновить подписку", http.StatusInternalServerError)
		return
	}

	// Получить обновленную подписку
	updated, err := h.service.GetSubscription(uint(id))
	if err != nil {
		http.Error(w, "Не удалось получить обновленную подписку", http.StatusInternalServerError)
		return
	}

	// Вернуть обновленную подписку
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DeleteSubscription удаляет подписку по ID
func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	// Получить ID из параметров URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID подписки", http.StatusBadRequest)
		return
	}

	// Удалить подписку из базы данных
	if err := h.service.DeleteSubscription(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Подписка не найдена", http.StatusNotFound)
			return
		}
		http.Error(w, "Не удалось удалить подписку", http.StatusInternalServerError)
		return
	}

	// Вернуть пустой ответ
	w.WriteHeader(http.StatusNoContent)
}

// ListSubscriptions получает список подписок с опциональными фильтрами и пагинацией
func (h *SubscriptionHandler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	// Получить параметры запроса
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	userID := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")

	// Установить значения по умолчанию
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// Получить список подписок
	subscriptions, total, err := h.service.ListSubscriptions(page, limit, userID, serviceName)
	if err != nil {
		http.Error(w, "Не удалось получить список подписок", http.StatusInternalServerError)
		return
	}

	// Вычислить общее количество страниц
	pages := int((total + int64(limit) - 1) / int64(limit))

	// Подготовить ответ
	response := struct {
		Data       []models.Subscription `json:"data"`
		Pagination struct {
			Page  int   `json:"page"`
			Limit int   `json:"limit"`
			Total int64 `json:"total"`
			Pages int   `json:"pages"`
		} `json:"pagination"`
	}{
		Data: subscriptions,
	}

	response.Pagination.Page = page
	response.Pagination.Limit = limit
	response.Pagination.Total = total
	response.Pagination.Pages = pages

	// Вернуть подписки
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CalculateTotalCost вычисляет общую стоимость подписок с опциональными фильтрами
func (h *SubscriptionHandler) CalculateTotalCost(w http.ResponseWriter, r *http.Request) {
	// Получить параметры запроса
	userID := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	// Вычислить общую стоимость
	totalCost, err := h.service.CalculateTotalCost(userID, serviceName, from, to)
	if err != nil {
		http.Error(w, "Не удалось вычислить общую стоимость: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Подготовить ответ
	response := struct {
		TotalCost int `json:"total_cost"`
	}{
		TotalCost: totalCost,
	}

	// Вернуть общую стоимость
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
