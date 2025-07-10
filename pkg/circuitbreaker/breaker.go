package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

//Шаблон Circuit_Breaker (Размыкатель цепи) автоматически отключает сервисные функции в ответ на вероятную неисправность,
//чтобы предотвратить более крупные или каскадные отказы, устранить повторяющиеся ошибки и обеспечить разумную реакцию на ошибки.

// Circuit — это функция, принимающая контекст и возвращающая строку или ошибку.
type Circuit func(context.Context) (string, error)

// Breaker возвращает обёртку для Circuit, реализующую механизм "предохранителя" (Circuit Breaker).
func Breaker(circuit Circuit, failureThreshold uint) Circuit {
	var consecutiveFailures int
	var lastAttempt = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (string, error) {
		m.RLock()
		d := consecutiveFailures - int(failureThreshold)
		if d >= 0 {
			// Рассчитываем время следующей разрешённой попытки
			shouldRetryAt := lastAttempt.Add(time.Second * time.Duration(1<<d))
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("service unreachable: circuit breaker is open")
			}
		}
		m.RUnlock()

		// Выполняем исходный Circuit
		response, err := circuit(ctx)

		// Переходим к записи
		m.Lock()
		defer m.Unlock()

		lastAttempt = time.Now()

		if err != nil {
			consecutiveFailures++
			return response, err
		}

		// Успешный ответ — сбрасываем счётчик
		consecutiveFailures = 0
		return response, nil
	}
}

//func sum(a, b int) (int, error) {
//	// Вдруг начинает ломаться
//	return 0, errors.New("произошла ошибка")
//}
//
//func main() {
//	// Адаптируем sum(a, b) под Circuit
//	circuit := func(ctx context.Context) (string, error) {
//		result, err := sum(3, 4)
//		if err != nil {
//			return "", err
//		}
//		return fmt.Sprintf("%d", result), nil
//	}
//
//	// Оборачиваем через Breaker с порогом 2 ошибки
//	breaker := Breaker(circuit, 2)
//
//	// Пытаемся вызвать 10 раз
//	for i := 1; i <= 10; i++ {
//		resp, err := breaker(context.Background())
//		fmt.Printf("Попытка %d: ответ = %q, ошибка = %v\n", i, resp, err)
//		time.Sleep(500 * time.Millisecond)
//	}
//}
