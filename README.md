# 🛡️ Go Simple Router with Recovery

## Library Go my custom router

Библиотека представляет собой минималистичный HTTP-роутер, реализованный поверх стандартной библиотеки Go.

Вместо использования готовых решений проект демонстрирует собственную реализацию механизма маршрутизации.

> [!NOTE]
> 
> ![Go Version](https://img.shields.io/badge/go-1.25%2B-blue.svg)
> ![License](https://img.shields.io/badge/license-GNU%20AGPLv3-red.svg)
> ![status: dev](https://img.shields.io/badge/status-dev-orange)
> ![CI](https://img.shields.io/badge/CI-GitHub%20Actions-green)
> ![CI](https://github.com/my-app-s/go-generator/actions/workflows/deploy.yml/badge.svg)
> ![REST API](https://img.shields.io/badge/REST%20API-green)

## TODO

На заметку unit-test

```text
Очень полезная привычка
Когда пишешь тест, не начинай с кода.
Сначала буквально напиши человеческое предложение:
«Если я создаю новый CORS, любой origin должен быть запрещён».
Потом превращаешь его:

создать CORS
       ↓
проверить origin
       ↓
ожидаем false

И только потом код.
Это сильно облегчает написание тестов.
```

- дать роутеру меч и щит
- проработать защиту
    - Recovery Middleware (Броня от падений)
		- Авторизация (JWT / Токены или Сессии)
		- Безопасность туннеля (Cloudflare)
		- Обязательно нужна защита от curl запросов, Postman или скриптов пример на Python и т.д.
    - Logger Middleware (Глаза и уши) Реализлвано но требует доработки
    - Rate Limiter (Щит от DDoS и брутфорса)
    - CORS Middleware (Пропуск для твоего фронта)
	- CORS не заменяет ***authentication*** и ***authorization*** обязательно!
		- Authentication
		- Authorization
		- Validation
		- Rate limiting

- Готово (для commit структуры)
	- Гибкость запростов по разным методам
		- structure map [string]map [string]
	- Обработка preflight-запросов
	- Логирование запросов
	- Реализована установка домена если нет то принимает от о Всех
		- Example
			```go
			r := router.NewRouterHandle() // create router
			r.AddRoute("/", welcomeHandler) // set path and handler
			r.Domain = "mydomain" // set domain
			wR := router.RecoveryMiddleware(r) // decorate Recovery
			wR = router.CorsMiddleware(wR, r.Domain) // decorate Cors
			```
			- обертка Recovery
			- обертка Cors
			- почищена ServeHTTP

	- Реализован черный ящик функция Handler структуры RouterHandle
		- В чем плюсы такого подхода:
			- Интерфейс «Черный ящик»: Тот, кто использует ваш пакет роутера, не должен задумываться о том, в каком порядке вешать Recovery и CORS, и что передавать в аргументы. Вся логика инкапсулирована внутри роутера.

			- Масштабируемость: Если завтра вам понадобится добавить логирование запросов (Logger), авторизацию или сжатие (Gzip), вы просто дописываете их внутрь метода r.Handler(), и вам не придется менять код во всех проектах, где используется этот роутер.

		- Example
		```go
		func main() {
			// 1. Создаем и настраиваем роутер
			r := router.NewRouterHandle()
			r.Domain = "mydomain"
			r.AddRoute("/", welcomeHandler)
			r.AddRoute("/about", aboutHandler)

			// 2. Одной строкой получаем полностью готовый, обернутый сервер
			handler := r.Handler()

			// 3. Запускаем
			log.Println("Server is running on :8080...")
			if err := http.ListenAndServe(":8080", handler); err != nil {
				log.Fatal(err)
			}
		}
		```

	- вынесено на слои разделение ответственности
		- Реализована защита RecoveryMiddleware обертка
			- добавить в инструкцию с предупреждением поднимать с RecoveryMiddleware
			- обязательна для оборачивания если не используется свой балансировщик выше
		- почищен ServeHTTP от defer для гибкости
		- example
			- ```go
				package main

				import (
					"log"
					"net/http"
					"your_project_name/router" // Замени на путь из твоего go.mod
				)

				func main() {
					// Создаем наш роутер
					r := router.NewRouter()

					// Добавляем тестовый маршрут
					r.Routes["/danger"] = func(w http.ResponseWriter, req *http.Request) {
						var ptr *string
						_ = *ptr // Паника! Сервер не упадет благодаря middleware
					}

					// Оборачиваем весь роутер в Recovery Middleware
					var protectedServer http.Handler = r
					protectedServer = router.RecoveryMiddleware(handler)

					log.Println("Сервер запущен на порту :8080...")
					if err := http.ListenAndServe(":8080", protectedServer); err != nil {
						log.Fatal(err)
					}
				}
				```

### Example for dev

```go
package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	r := chi.NewRouter()

	// 1. Щит от падений (Recovery)
	r.Use(middleware.Recoverer)

	// 2. Логгер запросов
	r.Use(middleware.Logger)

	// 3. Защита от зависаний (таймаут запроса)
	r.Use(middleware.Timeout(30 * time.Second))

	// 4. CORS (разрешаем запросы только с твоего GitHub Pages)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://твой-аккаунт.github.io"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Твои защищенные ручки
	r.Get("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Привет, это защищенный API!"))
	})

	http.ListenAndServe(":8080", r)
}
```

## 🚀 Основные возможности

-   **Panic Recovery**\
    Сервер не падает при ошибках --- используется `defer` + `recover`.

-   **O(1) Routing**\
    Мгновенный поиск маршрутов через `map`.

-   **Built-in Diagnostics**\
    Встроенный `/crashtest` для проверки отказоустойчивости.

-   **Fluent API**\
    Удобное добавление маршрутов через `AddRoute`.

-   **CORS**\
    Realization CORS request.

## 🛠️ Архитектура

Роутер реализует интерфейс `http.Handler`, что позволяет использовать
его напрямую в `http.ListenAndServe`.

### Как работает `ServeHTTP`:

1.  Инициализация `defer` + `recover`
2.  Поиск маршрута в `map`
3.  Вызов обработчика
4.  Если маршрут не найден → `404 Not Found`

## 💻 Использование

### Main example

``` go
package main

import (
    "fmt"
    "net/http"
    
    "github.com/my-app-s/go-custom-router/router"
)

func main() {
    r := router.NewRouterHandle()

    // Registering routes externally
    r.AddRoute("/", func(w http.ResponseWriter, req *http.Request) {
        w.Write([]byte("Welcome to custom router!\n"))
    })

    http.ListenAndServe(":8080", r)
}
```

### Flexible Route Registration

You can register routes using a standard call or a **Fluent API** (method chaining). Choose the style that fits your project best:

```go
r := router.NewRouterHandle()

// Standard style
r.AddRoute("/main", MainHandler)

// Fluent API style (Chaining)
r.
    AddRoute("/test", TestHandler).
    AddRoute("/login", LoginHandler).
    AddRoute("/profile", ProfileHandler)
```

## 🧪 Тестирование

``` bash
go test -v
go test -cover
```

### Как запустить BenchmarkRouter:

В терминале из папки с проектом выполните команду:

```bash
go test -bench=. -benchmem
```

> [!NOTE]
> Этот флаг -benchmem покажет вам две важнейшие вещи:
> - Сколько наносекунд занимает один запрос (ns/op).
> - Сколько байт памяти выделяется на один запрос и сколько было аллокаций (B/op и allocs/op).

### Handler signature
Your handlers should match the `http.HandlerFunc` signature:
```go
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, world!"))
}
```

## Disclaimer & License

- **Short Disclaimer (EN)**: Materials are provided ***as is*** under the LICENSE file. No warranties, no rights granted unless explicitly stated. Authors are not liable for damages. No partnership or obligations created.
- **Short Disclaimer (RU)**: Материалы предоставляются ***как есть*** и регулируются LICENSE. Гарантий нет, права не передаются без явного указания. Автор(ы) не несут ответственности. Партнёрство или обязательства не создаются.
- **Full Disclaimer**: Read the full text in the [DISCLAIMER.md](https://github.com/my-app-s/my-app-s/blob/main/DISCLAIMER.md) (Available in EN/RU).
- **License**: Distributed under the [GNU AGPLv3](./LICENSE) license.

## Author & Contacts

- **GitHub**: [@my-app-s](https://github.com/my-app-s)
- **LinkedIn**: [In/my-app-s](https://www.linkedin.com/in/my-app-s)
