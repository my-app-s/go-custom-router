# 🛡️ Go Simple Router with Recovery

## Library Go my custom router

Библиотека представляет собой минималистичный HTTP-роутер, реализованный поверх стандартной библиотеки Go.

Вместо использования готовых решений проект демонстрирует собственную реализацию механизма маршрутизации.

> [!NOTE]
> 
> ![Go Version](https://img.shields.io/badge/go-1.25%2B-blue.svg)
> ![License](https://img.shields.io/badge/license-GNU%20AGPLv3-red.svg)
> ![status: dev](https://img.shields.io/badge/status-dev-orange)

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

------------------------------------------------------------------------

## 📜 Disclaimer

Materials are provided **as is** under the terms of the LICENSE file. No warranties are provided. The authors are not liable for any damages.

📌 See the full disclaimer in [DISCLAIMER.md](./DISCLAIMER.md).

## 📜 License

Licensed under the GNU AGPLv3. See the [LICENSE](./LICENSE) file for details.

### Handler signature
Your handlers should match the `http.HandlerFunc` signature:
```go
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, world!"))
}
