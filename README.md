# 🛡️ Go Custom Router

Библиотека представляет собой легковесный, высокопроизводительный HTTP-роутер и набор утилит, реализованных поверх стандартной библиотеки Go (`net/http`), без использования сторонних тяжеловесных фреймворков.

---

> [!NOTE]
> 
> ![Educational/Experimental Project](https://img.shields.io/badge/Project-Educational%2FExperimental-green)
> ![Go Version](https://img.shields.io/badge/go-1.25%2B-blue.svg)
> ![License](https://img.shields.io/badge/license-GNU%20AGPLv3-red.svg)
> ![status: final](https://img.shields.io/badge/status-final-success)
> ![CI](https://img.shields.io/badge/CI-GitHub%20Actions-green)
> ![REST API](https://img.shields.io/badge/REST%20API-green)
> ![Latest Tag](https://img.shields.io/github/v/tag/my-app-s/go-custom-router)

---

## 🚀 Основные возможности

* **$O(1)$ Average-Case Routing**: Поиск маршрутов и методов основан на встроенных хэш-таблицах Go (`map`), что обеспечивает максимальную скорость диспетчеризации.
* **Fluent API**: Удобная цепочка вызовов (method chaining) для быстрой регистрации маршрутов (`GET`, `POST`, `PUT`, `DELETE`, `PATCH`).
* **Встроенный Rate Limiter**: Защита от DDoS и брутфорса на базе алгоритма *Sliding Window Counter* с автоматической очисткой старых записей воркером и поддержкой заголовков прокси (`CF-Connecting-IP`, `X-Forwarded-For`).
* **Middleware-пакет**: 
  - `RecoveryMiddleware` — перехват паник с логированием стека вызовов и возвратом `500 Internal Server Error`.
  - `LoggerMiddleware` — структурированное логирование времени выполнения и статусов запросов.
  - `ContentTypeJSONMiddleware` — автоматическая установка заголовка `application/json`.
* **Гибкий CORS**: Настройка разрешенных источников, методов, заголовков и предварительных запросов (`OPTIONS`).
* **JSON Helpers**: Удобные функции `SendJSON`, `SendError` и `MakeCustomHandler` для быстрой генерации ответов.

---

## 🛠️ Архитектура

Роутер полностью реализует стандартный интерфейс `http.Handler`, что позволяет бесшовно интегрировать его с любыми стандартными инструментами, балансировщиками и функциями тестирования (`net/http/httptest`).

### Порядок обработки запроса:
1. Логирование входящего запроса (`Logger`).
2. Перехват возможных паник (`Recovery`).
3. Установка заголовков (`CORS`, `Content-Type`).
4. Поиск совпадения по методу и пути в роутере.
5. Вызов целевого обработчика или возврат `404 Not Found` / `405 Method Not Allowed`.

---

## 💻 Примеры использования

### Базовый запуск с Fluent API

```go
package main

import (
	"log"
	"net/http"

	"[github.com/my-app-s/go-custom-router/router](https://github.com/my-app-s/go-custom-router/router)"
)

func main() {
	r := router.NewRouter()

	// Регистрация маршрутов через Fluent API
	r.
		GET("/", func(w http.ResponseWriter, req *http.Request) {
			_, _ = w.Write([]byte("Welcome to Go Custom Router!\n"))
		}).
		GET("/api/ping", router.MakeCustomHandler("status", "pong"))

	log.Println("Server is running on :8080...")
	if err := http.ListenAndServe(":8080", r.HandlerAPI()); err != nil {
		log.Fatal(err)
	}
}

```

### Использование Rate Limiter Middleware

```go
package main

import (
	"log"
	"net/http"
	"time"

	"[github.com/my-app-s/go-custom-router/router](https://github.com/my-app-s/go-custom-router/router)"
)

func main() {
	r := router.NewRouter()

	// Лимит: 100 запросов в минуту на один IP
	limiter := router.NewLimiter(100, 1*time.Minute, 5*time.Minute, 10000)
	defer limiter.Stop()

	r.GET("/api/limited", func(w http.ResponseWriter, req *http.Request) {
		router.SendJSON(w, http.StatusOK, map[string]string{"message": "Success"})
	})

	// Оборачиваем роутер или конкретный эндпоинт в лимитер
	handler := router.RateLimitMiddleware(limiter)(r)

	log.Println("Rate-limited server running on :8080...")
	_ = http.ListenAndServe(":8080", handler)
}

```

---

## 🧪 Тестирование

Проект покрыт модульными тестами с использованием встроенного пакета `httptest`. Для запуска тестов выполните:

```bash
go test -v ./...
go test -cover

```

---

## Disclaimer & License

* **Short Disclaimer (EN)**: Materials are provided ***as is*** under the LICENSE file. No warranties, no rights granted unless explicitly stated. Authors are not liable for damages. No partnership or obligations created.
* **Short Disclaimer (RU)**: Материалы предоставляются ***как есть*** и регулируются LICENSE. Гарантий нет, права не передаются без явного указания. Автор(ы) не несут ответственности. Партнёрство или обязательства не создаются.
* **Full Disclaimer**: Read the full text in the [DISCLAIMER.md](https://github.com/my-app-s/my-app-s/blob/main/DISCLAIMER.md) (Available in EN/RU).
* **License**: This project is dual-licensed:
- ​Open Source: [GNU AGPLv3](https://github.com/my-app-s/go-custom-router/blob/main/LICENSE)
- Commercial: A separate proprietary commercial license is available for proprietary and closed-source use. Contact the copyright holder for commercial licensing terms.

## Author & Contacts

* **GitHub**: [@my-app-s](https://github.com/my-app-s)
* **LinkedIn**: [In/my-app-s](https://www.linkedin.com/in/my-app-s)
