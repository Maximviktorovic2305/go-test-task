package middleware

import (
	"net/http"
	"runtime/debug"

	"effective-mobile-subscription/pkg/utils"
)

// ErrorMiddleware это промежуточное программное обеспечение для обработки ошибок и паник
func ErrorMiddleware(logger *utils.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Залогировать панику и трассировку стека
					logger.Error("Произошла паника",
						"ошибка", err,
						"стек", string(debug.Stack()),
						"url", r.URL.String(),
						"метод", r.Method,
					)

					// Отправить общий ответ об ошибке
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "Внутренняя ошибка сервера"}`))
				}
			}()

			// Обработать следующий обработчик
			next.ServeHTTP(w, r)
		})
	}
}
