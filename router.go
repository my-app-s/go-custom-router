// Copyright (C) 2025-2026 my-app-s
// Licensed under the GNU AGPLv3

// Package router provides a simple custom HTTP router
// for handling requests in Go applications.

package router

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// RouterHandle stores a mapping of URL paths to handler functions.
// It implements the http.Handler interface.
type RouterHandle struct {
	Routes map[string]http.HandlerFunc
}

// NewRouterHandle initializes a new RouterHandle with
// three default routes: "/", "/time", and "/date".
func NewRouterHandle() *RouterHandle {
	r := &RouterHandle{Routes: make(map[string]http.HandlerFunc)}
	r.Routes["/"] = r.HandleMain
	r.Routes["/time"] = r.HandleTime
	r.Routes["/date"] = r.HandleDate
	r.Routes["/crashtest"] = r.HandleCrashTest
	return r
}

// ServeHTTP dispatches incoming requests to the appropriate handler
// based on the request path. If no handler is found, it returns 404 Not Found.
func (r *RouterHandle) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Critical error: %v", err)
			http.Error(w, "Something broke on the server.", 500)
		}
	}()
	if handler, ok := r.Routes[req.URL.Path]; ok {
		handler(w, req)
	} else {
		http.NotFound(w, req)
	}
}

// HandleMain responds with basic request information including
// HTTP method, host, and path.
func (r *RouterHandle) HandleMain(res http.ResponseWriter, req *http.Request) {
	s := fmt.Sprintf("Method: %s\nHost: %s\nPath: %s",
		req.Method, req.Host, req.URL.Path)
	res.Write([]byte(s))
}

// HandleTime responds with the current system time
// formatted as DD.MM.YYYY HH:MM:SS.
func (r *RouterHandle) HandleTime(res http.ResponseWriter, req *http.Request) {
	s := time.Now().Format("02.01.2006 15:04:05")
	res.Write([]byte(s))
}

// HandleDate responds with the current system date
// formatted as DD.MM.YY.
func (r *RouterHandle) HandleDate(res http.ResponseWriter, req *http.Request) {
	s := time.Now().Format("02.01.06")
	res.Write([]byte(s))
}

// HandleCrashTest intentionally triggers a runtime panic
// by accessing an out-of-range index in a slice.
// Useful for testing router crash recovery and error handling.
func (r *RouterHandle) HandleCrashTest(res http.ResponseWriter, req *http.Request) {
	var list []int
	fmt.Println(list[99]) // exit to area list
}

// AddRoute allows you to add new routes and handlers after initialization
func (r *RouterHandle) AddRoute(path string, handler http.HandlerFunc) {
    r.Routes[path] = handler
}
