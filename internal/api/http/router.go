package http

import (
	"net/http"
)

func RegisterRoutes(h *BrokerHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /queue/{queueID}", h.Get)
	mux.HandleFunc("PUT /queue/{queueID}", h.Put)

	return mux
}
