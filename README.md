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

------------------------------------------------------------------------

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

------------------------------------------------------------------------

## 🛠️ Архитектура

Роутер реализует интерфейс `http.Handler`, что позволяет использовать
его напрямую в `http.ListenAndServe`.

### Как работает `ServeHTTP`:

1.  Инициализация `defer` + `recover`
2.  Поиск маршрута в `map`
3.  Вызов обработчика
4.  Если маршрут не найден → `404 Not Found`

------------------------------------------------------------------------

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

------------------------------------------------------------------------

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

------------------------------------------------------------------------

## Disclaimer & License

- **Short Disclaimer (EN)**: Materials are provided ***as is*** under the LICENSE file. No warranties, no rights granted unless explicitly stated. Authors are not liable for damages. No partnership or obligations created.
- **Short Disclaimer (RU)**: Материалы предоставляются ***как есть*** и регулируются LICENSE. Гарантий нет, права не передаются без явного указания. Автор(ы) не несут ответственности. Партнёрство или обязательства не создаются.
- **Full Disclaimer**: Read the full text in the [DISCLAIMER.md](https://github.com/my-app-s/my-app-s/blob/main/DISCLAIMER.md) (Available in EN/RU).
- **License**: Distributed under the [GNU AGPLv3](./LICENSE) license.

## Author & Contacts

- **GitHub**: [@my-app-s](https://github.com/my-app-s)
- **LinkedIn**: [In/my-app-s](https://www.linkedin.com/in/my-app-s)
