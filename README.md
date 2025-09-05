# Subscription Management Service

## Описание
Сервис для управления подписками пользователей с возможностями CRUDL и расчета стоимости.

## Технологии
- Go 1.25
- PostgreSQL
- GORM ORM
- Gorilla Mux
- Swagger UI

## Запуск

### Локально
1. Установить зависимости:
   ```
   go mod tidy
   ```

2. Настроить .env

3. Запустить приложение:
   ```
   go run cmd/server/main.go
   ```

### Docker
```
docker-compose up --build
```

## Документация API
- Swagger UI: http://localhost:8081/swagger/index.html
- Health check: http://localhost:8081/health