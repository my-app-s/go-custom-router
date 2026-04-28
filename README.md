# 🛡️ Go Simple Router with Recovery

Легкий, быстрый и устойчивый HTTP-роутер для приложений на Go.\
Проект ориентирован на простоту, надежность и удобство расширения.

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

## TODO

-   Add set method request []

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

## 📦 Стандартные маршруты

  Путь           Описание               Формат ответа
  -------------- ---------------------- -----------------------
  `/`            Информация о запросе   Method, Host, Path
  `/time`        Текущее время          `DD.MM.YYYY HH:MM:SS`
  `/date`        Текущая дата           `DD.MM.YY`
  `/crashtest`   Эмуляция аварии        Вызывает `panic`

------------------------------------------------------------------------

## 💻 Использование

### Main example

``` go
package main

import (
    "net/http"
    "github.com/my-app-s/go-custom-router/router"
)

func main() {
    r := router.NewRouterHandle()

    r.AddRoute("/main", MainHandler)

    println("Server is running on :8080")
    http.ListenAndServe(":8080", r)
}
```

### Flexible Route Registration
You can register routes using a standard call or a **Fluent API** (method chaining). 
Choose the style that fits your project best:

``` go
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

------------------------------------------------------------------------

## ⚠️ Disclaimer / Отказ от ответственности

### English Version
This project is an **independent development** provided on an **"AS IS"** basis.

* **Liability:** In no event shall the author be liable for any errors, bugs, or data loss arising from the use of this software.
* **Status:** This is an experimental tool. Always verify the generated HTML output before deployment.

> [!CAUTION]
> Any use (operation) of this code is at your own risk.

---

### Русская версия
Данный проект является **независимой разработкой** и предоставляется «как есть».

* **Ответственность:** Автор не несет ответственности за любые ошибки, баги или потерю данных, возникшие в результате использования данного кода.
* **Статус:** Это экспериментальный инструмент. Проверяйте сгенерированный HTML перед деплоем.

> [!CAUTION]
> Любое использование (эксплуатация) данного кода осуществляется на ваш страх и риск.

------------------------------------------------------------------------

## 📜 License

GNU AGPLv3

### Handler signature
Your handlers should match the `http.HandlerFunc` signature:
```go
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, world!"))
}
