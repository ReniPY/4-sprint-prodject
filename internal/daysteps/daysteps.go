package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

func parsePackage(data string) (int, time.Duration, error) {
	// ваш код ниже
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("ожидалось два элемента, но получили %d", len(parts))
	}

	// Преобразуем первый элемент в int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("не удалось преобразовать количество шагов")
	}

	// Преобразуем второй элемент в time.Duration
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("не удалось преобразовать продолжительность")
	}

	return steps, duration, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	// Разбираем данные
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("Ошибка при разборе данных:", err)
		return ""
	}

	// Проверяем, чтобы количество шагов было больше 0
	if steps <= 0 {
		fmt.Println("Количество шагов должно быть больше 0")
		return ""
	}

	// Вычисляем дистанцию в метрах
	distance := float64(steps) * StepLength

	// Переводим дистанцию в километры
	distanceKm := distance / 1000

	// Вычисляем потраченные калории
	roundedWeight := int(weight + 0.5)
	calories := spentcalories.WalkingSpentCalories(steps, float64(roundedWeight), height, duration)

	// Формируем строку с результатами
	result := fmt.Sprintf("Количество шагов: %d.\n Дистанция составила %.2f км.\n Вы сожгли %.2f калорий.", steps, distanceKm, calories)

	return result
}
