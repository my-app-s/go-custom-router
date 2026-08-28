# 🛡️ Go Custom Router

The library is a lightweight, high-performance HTTP router and set of utilities implemented on top of the standard Go library (`net/http`), without the use of third-party heavyweight frameworks.

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

## 🚀 Main features

* **$O(1)$ Average-Case Routing**: Route and method lookup is based on Go's built-in hash tables (`map`) for maximum dispatch speed.
* **Fluent API**: Convenient method chaining for quick registration of routes (`GET`, `POST`, `PUT`, `DELETE`, `PATCH`).
* **Built-in Rate Limiter**: Protection against DDoS and brute force based on the *Sliding Window Counter* algorithm with automatic cleaning of old records by the worker and support for proxy headers (`CF-Connecting-IP`, `X-Forwarded-For`).
* **Middleware package**: 
- `RecoveryMiddleware` - intercepting panics with call stack logging and returning `500 Internal Server Error`. 
- `LoggerMiddleware` - structured logging of execution time and request statuses. 
- `ContentTypeJSONMiddleware` - automatic setting of the `application/json` header.
* **Flexible CORS**: Configure allowed origins, methods, headers and pre-requests (`OPTIONS`).
* **JSON Helpers**: Convenient functions `SendJSON`, `SendError` and `MakeCustomHandler` for quickly generating responses.

---

## 🛠️ Architecture

The router fully implements the standard `http.Handler` interface, which allows it to be seamlessly integrated with any standard tools, balancers and testing functions (`net/http/httptest`).

### Request processing order:
1. Logging an incoming request (`Logger`).
2. Intercepting possible panics (`Recovery`).
3. Setting headers (`CORS`, `Content-Type`).
4. Search for a match by method and path in the router.
5. Calling the target handler or returning `404 Not Found` / `405 Method Not Allowed`.

---

## 💻 Examples of use

### Basic launch with Fluent API

```go
package main

import (
	"log"
	"net/http"

	"github.com/my-app-s/go-custom-router/router"
)

func main() {
	r := router.NewRouter()

	// Register routes via Fluent API
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

### Using Rate Limiter Middleware

```go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/my-app-s/go-custom-router/router"
)

func main() {
	r := router.NewRouter()

	// Limit: 100 requests per minute per IP
	limiter := router.NewLimiter(100, 1*time.Minute, 5*time.Minute, 10000)
	defer limiter.Stop()

	r.GET("/api/limited", func(w http.ResponseWriter, req *http.Request) {
		router.SendJSON(w, http.StatusOK, map[string]string{"message": "Success"})
	})

	// Wrap the router or a specific endpoint in the limiter
	handler := router.RateLimitMiddleware(limiter)(r)

	log.Println("Rate-limited server running on :8080...")
	_ = http.ListenAndServe(":8080", handler)
}

```

---

## 🧪 Testing

The project is covered with unit tests using the built-in `httptest` package. To run tests run:

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
* **Mail**: [myapps.mre.dev@gmail.com](mailto:myapps.mre.dev@gmail.com)
