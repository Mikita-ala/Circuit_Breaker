package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

func main() {
	var i int
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		if rand.Float32() < 0.5 {
			fmt.Printf("[%d] Ответ: %q\n", i, "internal error")
			i++
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		fmt.Printf("[%d] Ответ: %q\n", i, "Hello, world!")
		i++
		fmt.Fprintln(w, "Hello, world!")
	})

	fmt.Println("Сервер запущен на :8081")
	http.ListenAndServe(":8081", nil)
}
