package http

import (
	"encoding/xml"
	"log/slog"
	"net/http"

	"worker/internal/domain"
)

type Router struct {
	taskCrackerUseCase domain.ITaskCrackerUseCase
	logger             *slog.Logger
}

func NewRouter(taskCrackerUseCase domain.ITaskCrackerUseCase, logger *slog.Logger) *Router {
	return &Router{
		taskCrackerUseCase: taskCrackerUseCase,
		logger:             logger,
	}
}

func (r *Router) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /internal/api/worker/hash/crack/task", r.crackTask)

	return mux
}

func (r *Router) crackTask(w http.ResponseWriter, req *http.Request) {
	var workerRequest WorkerRequest
	if err := xml.NewDecoder(req.Body).Decode(&workerRequest); err != nil {
		r.logger.Error("failed to decode worker request", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	task := domain.Task{
		RequestId:  workerRequest.RequestId,
		Hash:       workerRequest.Hash,
		Alphabet:   workerRequest.Alphabet,
		MaxLength:  workerRequest.MaxLength,
		PartNumber: workerRequest.PartNumber,
		PartCount:  workerRequest.PartCount,
	}

	if err := r.taskCrackerUseCase.CrackTask(task); err != nil {
		r.logger.Error("failed to accept task", "error", err)
		http.Error(w, "failed to accept task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
