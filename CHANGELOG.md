## [v1.4.2] - 2026-08-26

### Added
- **Пакет `router` (CORS)**:
  - Структура `CORS` для управления политиками междоменного доступа.
  - Конструктор `NewCORS()` с гарантированной инициализацией карты разрешенных доменов.
  - Метод `AddOrigin()` для регистрации доменов с валидацией пустых значений.
  - Безопасный метод `IsAllowed()` с защитой от `nil`-указателей.
  - `CorsMiddleware()` — HTTP-middleware для выборочной фильтрации Origin и автоматической обработки preflight-запросов (`OPTIONS`).
  - `CorsMiddlewareOpen()` — публичный HTTP-middleware для открытых API (`Origin: *`).
  - Полная GoDoc-документация для всех экспортируемых типов и функций.
- **Тесты (`router/cors_test.go`)**:
  - Табличные модульные тесты (Table-driven tests) для проверки валидации и добавления доменов (`TestAddOrigin`).
  - Тест инициализации конструктора (`TestNewCORS`).
  - Проверка корректности работы с дубликатами (`TestMultiAddOrigin`).
  - Интеграционные HTTP-тесты с использованием `httptest.NewRecorder` для проверки заголовков ответа и метода `OPTIONS` (`TestCorsMiddleware`, `TestCorsMiddlewareOpen`).

## [1.5.2] - 2026-08-26

### Added
- **Модуль Middleware (`router/middleware.go`)**:
  - `RecoveryMiddleware()` — перехват паник (`panic`) с логированием stack trace и возвратом 500 Internal Server Error.
  - `LoggerMiddleware()` — логирование HTTP-методов, путей и времени выполнения запроса (`[START]` / `[DONE]`).
  - `ContentTypeJSONMiddleware()` — автоматическая установка заголовка `Content-Type: application/json; charset=utf-8`.
  - Методы `Handler()` и `HandlerAPI()` для автоматической сборки цепочек middleware внутри `RouterHandle`.
  - Полная GoDoc-документация для всех экспортируемых функций и методов.
- **Тестирование (`router/middleware_test.go`)**:
  - Тесты изоляции паник (`TestRecoveryMiddleware`).
  - Проверка логирования и заголовков (`TestLoggerMiddleware`, `TestContentTypeJSONMiddleware`).
  - Тесты сборщиков цепочек (`TestRouterHandle_Handler`, `TestRouterHandle_HandlerAPI`).

### Changed
- Обновлены заголовки авторских прав (Copyright 2025-2026) и указания лицензии GNU AGPLv3 в исходных файлах.

## [1.5.3] - 2026-08-26

### Added
- **Маршрутизация (`router/routes.go`)**:
  - Метод `AddRoute()` и обёртки-помощники `GET()`, `POST()`, `PUT()`, `DELETE()` с поддержкой Fluent API.
  - Набор модульных тестов для проверки карты маршрутов (`router/routes_test.go`).
  - Добавлены GoDoc-комментарии и AGPLv3 заголовки лицензии.

## [1.5.4] - 2026-08-26

### Added
- **Ядро роутера (`router/router.go`)**:
  - В структуру `RouterHandle` добавлено поле `CORS *CORS` для связки с подсистемой middleware.
  - Конструктор `NewRouterHandle()` с инициализацией карты маршрутов.
  - Диспетчер `ServeHTTP` с корректной обработкой статусов `404 Not Found` и `405 Method Not Allowed`.
  - Модульные тесты маршрутизации и кодов ответов (`router/router_test.go`).
