package handlers

import "net/http"

type Handler interface {
	http.Handler
	Route() string
	Method() []string
}
