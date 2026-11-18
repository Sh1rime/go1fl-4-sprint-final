package spentcalories

import (
	"time"
	"strings"
	"fmt"
	"strconv"
	"log"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid data format: expected 3 fields, got %d", len(parts))
	}
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("steps count must be greater than zero")
	}
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("time must be greater than zero")
	}
	return steps, parts[1], duration, nil
}

func distance(steps int, height float64) float64 {
	return height * stepLengthCoefficient * float64(steps) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, typeTraining, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	distance := distance(steps, height)
	if distance <= 0 {
		return "", fmt.Errorf("less than or equal to zero: distance %.2f", distance)
	}
	meanSpeed := meanSpeed(steps, height, duration)
	if meanSpeed <= 0 {
		return "", fmt.Errorf("less than or equal to zero: meanSpeed %.2f", meanSpeed)
	}
	var calories float64
	switch typeTraining {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",typeTraining, duration.Hours(), distance, meanSpeed, calories), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("value less than or equal to zero: steps %d, weight %.2f, height %.2f, duration %v", steps, weight, height, duration)
	}
	return (weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH, nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("value less than or equal to zero: steps %d, weight %.2f, height %.2f, duration %v", steps, weight, height, duration)
	}
	return (weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH * walkingCaloriesCoefficient, nil
}
