# ⚡ Circuit Breaker (Размыкатель цепи) на Go

Проект демонстрирует применение шаблона проектирования **Circuit Breaker** на языке Go. Механизм автоматически отключает сервисные функции после нескольких неудачных вызовов, чтобы избежать каскадных сбоев и чрезмерной нагрузки на зависимый сервис.

---

## 📁 Структура проекта

```
Circuit_Breaker/
├── cmd/
│   ├── client/         # HTTP-клиент, обёрнутый в Circuit Breaker
│   │   └── main.go
│   └── server/         # Нестабильный HTTP-сервер
│       └── main.go
├── pkg/
│   └── circuitbreaker/ # Реализация Circuit Breaker
│       └── breaker.go
├── go.mod              # Модуль Go
└── README.md
```

---

## 🔧 Установка и запуск

Открой два терминала.

### 1. Запусти сервер:

```bash
go run ./cmd/server
```

Сервер будет слушать на `localhost:8081/hello`. Он имитирует сбои — с вероятностью 50% возвращает `500 Internal Server Error`.

### 2. Запусти клиент:

```bash
go run ./cmd/client
```

Клиент будет запрашивать сервер каждые 1 секунду, используя обёртку Circuit Breaker. Если количество подряд неудачных вызовов превышает порог, дальнейшие вызовы временно блокируются.

---

## 🛠 Механизм Circuit Breaker

Модуль `pkg/circuitbreaker` реализует простую обёртку `Breaker`, принимающую функцию `Circuit`:

```go
type Circuit func(ctx context.Context) (string, error)

func Breaker(circuit Circuit, failureThreshold uint) Circuit
```

Если количество **последовательных неудач** ≥ `failureThreshold`, Circuit Breaker **отключает** вызовы и начинает возвращать:

```text
"service unreachable: circuit breaker is open"
```

Через определённое **экспоненциальное время** (2, 4, 8, 16 сек и т.д.) он пробует снова. Успешный вызов сбрасывает счётчик.

---

## 🖥 Пример работы клиента:

```text
[0] Ответ: "", Ошибка: unexpected status code: 500
[1] Ответ: "Hello, world!\n"
...
[27] Ответ: "", Ошибка: unexpected status code: 500
[28] Ответ: "", Ошибка: service unreachable: circuit breaker is open
[29] Ответ: "Hello, world!\n"
```

## 🛰 Пример вывода сервера:

```text
[0] Ответ: "internal error"
[1] Ответ: "Hello, world!"
...
[28] Ответ: "Hello, world!"
```

---

## 📌 Полезно знать

- Circuit Breaker особенно полезен в микросервисной архитектуре, где отказ одной службы может повлечь лавинообразный сбой всей системы.
- Реализация легко расширяется под другие вызовы: базы данных, gRPC, Redis и т. д.

---

## ✅ Пример использования в другом проекте:

```go
breaker := circuitbreaker.Breaker(myFunc, 3)

for {
    result, err := breaker(context.Background())
    if err != nil {
        log.Println("Ошибка:", err)
        continue
    }
    fmt.Println("Результат:", result)
}
```

---

## 📚 Источники

- [Cloud Native Go Matthew A. Titmus](https://ftp.zhirov.kz/books/IT/Go/Облачный%20Go%20%28М.А.%20Титмус%29.pdf) — оригинальное описание шаблона Circuit Breaker.
- [Go net/http Documentation](https://pkg.go.dev/net/http)

---
