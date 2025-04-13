package http

import (
	"encoding/json"
	"encoding/xml"
	"log/slog"
	"net/http"

	"manager/internal/domain"
)

type Router struct {
	hashUseCase domain.IHashUseCase
	logger      *slog.Logger
}

func NewRouter(hashUseCase domain.IHashUseCase, logger *slog.Logger) *Router {
	return &Router{
		hashUseCase: hashUseCase,
		logger:      logger,
	}
}

func (r *Router) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/hash/crack", r.crackHash)
	mux.HandleFunc("GET /api/hash/status", r.getStatus)
	mux.HandleFunc("PATCH /internal/api/manager/hash/crack/request", r.updateRequest)

	return mux
}

func (r *Router) crackHash(w http.ResponseWriter, req *http.Request) {
	var request CrackHashRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		r.logger.Error("failed to decode request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	hashData, err := r.hashUseCase.CrackHash(request.Hash, request.MaxLength)
	if err != nil {
		r.logger.Error("failed to process hash request", "error", err)
		http.Error(w, "failed to process request", http.StatusInternalServerError)
		return
	}

	response := CrackHashResponse{RequestId: hashData.RequestId}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		r.logger.Error("failed to encode response", "error", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (r *Router) getStatus(w http.ResponseWriter, req *http.Request) {
	requestId := req.URL.Query().Get("requestId")
	if requestId == "" {
		r.logger.Error("missing requestId parameter")
		http.Error(w, "missing requestId parameter", http.StatusBadRequest)
		return
	}

	hashData, err := r.hashUseCase.GetStatus(requestId)
	if err != nil {
		r.logger.Error("failed to get status", "error", err)
		http.Error(w, "failed to get status", http.StatusInternalServerError)
		return
	}

	response := CrackHashStatus{
		Status: string(hashData.Status),
	}
	if hashData.Status != domain.ErrorStatus {
		response.Data = hashData.Data
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		r.logger.Error("failed to encode response", "error", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (r *Router) updateRequest(w http.ResponseWriter, req *http.Request) {
	var workerResponse WorkerResponse
	if err := xml.NewDecoder(req.Body).Decode(&workerResponse); err != nil {
		r.logger.Error("failed to decode worker response", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := r.hashUseCase.UpdateRequest(workerResponse.RequestId, workerResponse.Words); err != nil {
		r.logger.Error("failed to update request status", "error", err)
		http.Error(w, "failed to update request status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
