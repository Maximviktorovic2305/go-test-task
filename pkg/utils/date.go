package utils

import (
	"fmt"
	"time"
)

// ParseMonthYear разбирает строку даты в формате ММ-ГГГГ
func ParseMonthYear(dateStr string) (time.Time, error) {
	// Разобрать строку даты в формате ММ-ГГГГ
	t, err := time.Parse("01-2006", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("недопустимый формат даты, ожидается ММ-ГГГГ: %v", err)
	}
	return t, nil
}

// FormatMonthYear форматирует time.Time в строку ММ-ГГГГ
func FormatMonthYear(t time.Time) string {
	return t.Format("01-2006")
}

// GetFirstDayOfMonth возвращает первый день месяца для заданной даты
func GetFirstDayOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// GetLastDayOfMonth возвращает последний день месяца для заданной даты
func GetLastDayOfMonth(t time.Time) time.Time {
	// Получить первый день следующего месяца
	nextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	// Вычесть один день, чтобы получить последний день текущего месяца
	return nextMonth.AddDate(0, 0, -1)
}
