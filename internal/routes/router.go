package routes

import (
	"net/http"

	"effective-mobile-subscription/internal/handlers"
	"effective-mobile-subscription/internal/repository"
	"effective-mobile-subscription/internal/services"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// SetupRoutes настраивает все маршруты для приложения
func SetupRoutes(db *gorm.DB) *mux.Router {
	// Создать маршрутизатор
	router := mux.NewRouter()

	// Создать репозитории
	subscriptionRepo := repository.NewSubscriptionRepository(db)

	// Создать сервисы
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	// Создать обработчики
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)

	// Настроить маршруты
	setupSubscriptionRoutes(router, subscriptionHandler)

	// Проверка состояния
	// swagger:operation GET /health health healthCheck
	//
	// Check the health of the service
	//
	// Returns a simple status message to indicate the service is running
	//
	// ---
	// responses:
	//   '200':
	//     description: Service is healthy
	//     schema:
	//       type: object
	//       properties:
	//         status:
	//           type: string
	//           example: ok
	router.HandleFunc("/health", healthCheck).Methods("GET")

	return router
}

// setupSubscriptionRoutes настраивает маршруты для управления подписками
func setupSubscriptionRoutes(router *mux.Router, handler *handlers.SubscriptionHandler) {
	// Операции CRUDL
	router.HandleFunc("/subscriptions", handler.CreateSubscription).Methods("POST")
	router.HandleFunc("/subscriptions/{id:[0-9]+}", handler.GetSubscription).Methods("GET")
	router.HandleFunc("/subscriptions/{id:[0-9]+}", handler.UpdateSubscription).Methods("PUT")
	router.HandleFunc("/subscriptions/{id:[0-9]+}", handler.DeleteSubscription).Methods("DELETE")
	router.HandleFunc("/subscriptions", handler.ListSubscriptions).Methods("GET")

	// Расчет стоимости
	router.HandleFunc("/subscriptions/cost", handler.CalculateTotalCost).Methods("GET")
}

// healthCheck - проверка состояния
//
//	@Summary		Health check
//	@Description	Returns a simple status message to indicate the service is running
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	object{status=string}
//	@Router			/health [get]
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}
