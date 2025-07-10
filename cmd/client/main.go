package main

import (
	"Circuit_Breaker/pkg/circuitbreaker"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{Timeout: 1 * time.Second}

	circuit := func(ctx context.Context) (string, error) {
		// Создаётся HTTP-запрос GET к http://localhost:8081/hello, с поддержкой отмены через context.Context.
		req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8081/hello", nil)
		if err != nil {
			return "", err
		}
		resp, err := client.Do(req) // Делаем сам HTTP-запрос через http.Client. Это может занять время — зависит от сервера.
		if err != nil {
			return "", err
		}
		defer resp.Body.Close() // Гарантируем, что тело ответа будет закрыто — даже при ошибке чтения.

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		// Читаем тело ответа.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		// Если всё прошло успешно — возвращаем ответ сервера как строку.
		return string(body), nil
	}

	breaker := circuitbreaker.Breaker(circuit, 2)
	//var i int
	//for {
	//	result, err := breaker(context.Background())
	//	fmt.Printf("[%d] Ответ: %q, Ошибка: %v\n", i, result, err)
	//	i++
	//	time.Sleep(1 * time.Second)
	//}

	for i := 0; i < 30; i++ {
		result, err := breaker(context.Background())
		if err != nil {
			fmt.Printf("[%d] Ответ: %q, Ошибка: %v\n", i, result, err)
		} else {
			fmt.Printf("[%d] Ответ: %q\n", i, result)
		}
		time.Sleep(1 * time.Second)
	}
}
