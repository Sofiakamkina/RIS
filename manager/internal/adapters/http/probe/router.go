package probe

import (
	"log/slog"
	"net/http"
)

const (
	healthURL = "/healthz"
)

type Router struct {
	logger *slog.Logger
}

func NewRouter(logger *slog.Logger) *Router {
	return &Router{
		logger: logger,
	}
}

func (r *Router) Health(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("success")); err != nil {
		r.logger.Error("failed to write probe response", "err", err)
	}
}

func (r *Router) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(healthURL, r.Health)

	return mux
}
