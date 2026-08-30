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

## [v1.6.0] - 2026-08-26

### Added
- **Защита от флуда и DDoS (`router/rate_limiter.go`)**:
  - `CustomRateLimiter` с алгоритмом временных окон и фоновой очисткой памяти.
  - Извлечение реального IP из заголовков `CF-Connecting-IP`, `X-Forwarded-For` и очистка портов `RemoteAddr`.
  - Middleware `RateLimitMiddleware()` с возвратом статуса `429 Too Many Requests`.
  - Модульные и HTTP-тесты логики ограничений (`router/rate_limiter_test.go`).

## [1.6.5] - 2026-08-26

### Added
- **Вспомогательные функции (`router/helpers.go`)**:
  - `SendJSON()` для кодирования и отправки JSON-ответов.
  - `SendError()` для стандартизированных JSON-ошибок.
  - `MakeCustomHandler()` для генерации mock-хэндлеров и статических ответов.
  - Модульные тесты хелперов (`router/helpers_test.go`).

## [1.6.6] - 2026-08-26
### Added
- Comprehensive architecture notes, TODOs, and developer examples in `README.md`
- Dependency `golang.org/x/time` for rate limiting features

### Fixed
- Removed duplicate `RouterHandle` struct declaration in `router/middleware.go`
- Fixed route registration and 404 errors in middleware tests by using `router.GET()`
- Stabilized full-chain middleware tests (`Recovery`, `CORS`, `Logger`, `API`)

## [1.6.7] - 2026-08-26

### Fixed
- Updated README structure and code blocks formatting

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [1.6.8] - 2026-08-28

### Added

* Added the new `Router` type as the primary HTTP router implementation.
* Added `NewRouter()` constructor for router initialization.
* Added the generic `Handle()` method for route registration using `http.Handler`.
* Added route helper methods for `GET`, `POST`, `PUT`, `DELETE`, and `PATCH`.
* Added normalized route handling using `path.Clean()`.
* Added per-origin CORS configuration through the new `OriginsOptions` type.
* Added configurable `RateLimiter` with:

  * maximum request count;
  * configurable time window;
  * cleanup TTL;
  * maximum number of tracked clients.
* Added sliding-window request estimation to the rate limiter.
* Added graceful shutdown support for the rate limiter cleanup worker through `Stop()`.
* Added `Retry-After` header to rate-limited HTTP 429 responses.
* Added validation of client IP addresses obtained from `CF-Connecting-IP`, `X-Forwarded-For`, and `RemoteAddr`.
* Added tests for `SendError()` JSON responses.
* Added additional CORS response header assertions.
* Added cleanup-worker shutdown handling to rate limiter tests.

### Changed

#### Router

* Replaced `RouterHandle` with `Router`.
* Changed internal route storage from the exported `Routes` field to an unexported `routes` map.
* Changed route handlers to use the standard `http.Handler` interface.
* Changed request dispatching to use `handler.ServeHTTP()`.
* Route registration and request lookup now normalize paths with `path.Clean()`.
* Updated the middleware chain to operate on the new `Router` implementation.

#### CORS

* Changed CORS configuration from global `AllowedMethods` and `AllowedHeaders` to per-origin configuration.
* `AllowedOrigins` now stores `OriginsOptions` for each configured origin.
* `AddOrigin()` now associates allowed methods and headers directly with the corresponding origin.
* Added automatic initialization of the `AllowedOrigins` map when required.
* Expanded open CORS support to `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, and `OPTIONS`.
* Added `Authorization` to the default allowed CORS headers.

#### Rate Limiting

* Replaced `CustomRateLimiter` with `RateLimiter`.
* Replaced `NewCustomRateLimiter()` with `NewLimiter()`.
* Added `MaxClients` to limit the number of tracked client entries.
* Changed the cleanup interval from 1 minute to 30 seconds.
* Added explicit ticker cleanup with `ticker.Stop()`.
* Reworked request limiting from a fixed-window counter to a sliding-window estimation approach using previous and current window counters.
* Added a dedicated stop channel for terminating the background cleanup worker.

#### Middleware

* Updated `HandlerAPI()` to use the new `Router` implementation.
* Added optional rate limiting to the API middleware chain when a `RateLimiter` is configured.
* Updated middleware tests for the new router API.
* Removed request logging statements from the current `LoggerMiddleware` implementation.
* Simplified middleware documentation comments.

#### Tests

* Migrated router tests from `RouterHandle` to `Router`.
* Updated rate limiter tests to use `NewLimiter()`.
* Improved JSON unmarshalling error handling in tests.
* Changed route tests to verify actual HTTP dispatch through `ServeHTTP()`.
* Updated CORS tests to validate per-origin method configuration.

### Removed

* Removed the legacy `RouterHandle` type.
* Removed `NewRouterHandle()`.
* Removed the legacy `AddRoute()` API.
* Removed `CustomRateLimiter`.
* Removed `NewCustomRateLimiter()`.
* Removed global CORS `AllowedMethods` and `AllowedHeaders` fields.
* Removed tests targeting the legacy `RouterHandle` API.

### Breaking Changes

> **Warning:** This release contains breaking API changes.

* `RouterHandle` has been replaced by `Router`.
* `NewRouterHandle()` has been replaced by `NewRouter()`.
* `AddRoute()` has been replaced by `Handle()`.
* The public `Routes` field is no longer available.
* `CustomRateLimiter` has been replaced by `RateLimiter`.
* `NewCustomRateLimiter()` has been replaced by `NewLimiter()`.
* `NewLimiter()` now requires an additional `maxClients` argument.
* CORS configuration has changed from global method/header settings to per-origin `OriginsOptions`.
* `HandlerAPI()` now optionally applies rate limiting when a limiter is configured.

### Documentation

* Updated routing documentation to describe map lookup as **O(1) average-case**.
* Added documentation for the sliding-window rate limiting approach.
* Documented the project's educational/experimental status.
* Updated the README licensing section to describe GNU AGPLv3 open-source licensing and a separate proprietary commercial licensing option.
  [1.6.8]: https://github.com/my-app-s/go-custom-router/compare/v1.6.7...v1.6.8
  [1.6.7]: https://github.com/my-app-s/go-custom-router/releases/tag/v1.6.7

## [1.6.10] - 2026-08-30

### Fixed
- Polished disclaimer wording and formatted license links in README
