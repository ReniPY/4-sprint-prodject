package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// Разделяем строку на три части
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("ожидалось три элемента, но получили %d", len(parts))
	}

	// Преобразуем первый элемент в int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось преобразовать количество шагов")
	}

	// Преобразуем третий элемент в time.Duration
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось преобразовать продолжительность")
	}

	// Вид активности
	activity := parts[1]

	return steps, activity, duration, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	// Вычисляем дистанцию в метрах
	distanceInMeters := float64(steps) * lenStep
	// Переводим в километры
	return distanceInMeters / mInKm
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration == 0 {
		return 0
	}
	// Вычисляем дистанцию
	distance := distance(steps)
	// Переводим продолжительность в часы
	hours := duration.Hours()
	// Вычисляем среднюю скорость
	return distance / hours
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	// Парсим данные
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return fmt.Sprintf("Ошибка при разборе данных: %v", err)
	}

	// Вычисляем дистанцию
	distance := distance(steps)
	// Вычисляем среднюю скорость
	speed := meanSpeed(steps, duration)

	var calories float64
	switch activity {
	case "Бег":
		calories = RunningSpentCalories(steps, weight, duration)
	case "Ходьба":
		calories = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "неизвестный тип тренировки"
	}

	// Формируем строку с результатами
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f", activity, duration.Hours(), distance, speed, calories)
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	// Вычисляем среднюю скорость
	meanSpeed := meanSpeed(steps, duration)
	// Вычисляем калории
	calories := ((runningCaloriesMeanSpeedMultiplier * meanSpeed) - runningCaloriesMeanSpeedShift) * weight
	return calories
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	// Вычисляем среднюю скорость
	speed := meanSpeed(steps, duration)
	// Переводим продолжительность в часы
	hours := duration.Hours()
	// Вычисляем калории
	calories := ((walkingCaloriesWeightMultiplier * weight) + (speed*speed/height)*walkingSpeedHeightMultiplier) * hours * minInH
	return calories
}
