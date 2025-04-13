package domain

import (
	"errors"
	"log/slog"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"
)

type HashUseCase struct {
	mutex        sync.Mutex
	requests     map[string]HashData
	workersCount int
	workersURL   []string
	alphabet     string
	ttl          int
	repository   IRepository
	logger       *slog.Logger
}

func NewHashUseCase(
	workersCount int,
	workersURL []string,
	alphabet string,
	ttl int,
	repository IRepository,
	logger *slog.Logger,
) *HashUseCase {
	return &HashUseCase{
		requests:     make(map[string]HashData),
		workersCount: workersCount,
		workersURL:   workersURL,
		alphabet:     alphabet,
		ttl:          ttl,
		repository:   repository,
		logger:       logger,
	}
}

func (h *HashUseCase) CrackHash(hash string, maxLength int) (HashData, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	requestId := uuid.New().String()

	hashData := HashData{
		RequestId: requestId,
		Status:    InProgressStatus,
	}

	h.requests[requestId] = hashData

	tasks := h.splitTask(requestId, hash, maxLength)

	successChan := make(chan bool, h.workersCount)

	for i, task := range tasks {
		go func(i int, task WorkerTask) {
			if err := h.repository.SendTaskToWorker(h.workersURL[i], task); err != nil {
				h.logger.Error("failed to send task to worker", slog.String("workerURL", h.workersURL[i]), slog.String("error", err.Error()))
				for j := 0; j < h.workersCount; j++ {
					if j == i {
						continue
					}
					if err = h.repository.SendTaskToWorker(h.workersURL[j], task); err == nil {
						successChan <- true
						return
					}
				}
				successChan <- false
			} else {
				successChan <- true
			}
		}(i, task)
	}

	for i := 0; i < h.workersCount; i++ {
		if !<-successChan {
			return hashData, errors.New("all workers are unavailable")
		}
	}

	go h.startRequestTimeout(requestId)

	return hashData, nil
}

func (h *HashUseCase) GetStatus(requestId string) (HashData, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	response, ok := h.requests[requestId]
	if !ok {
		return HashData{}, errors.New("request not found")
	}

	return response, nil
}

func (h *HashUseCase) UpdateRequest(requestId string, data []string) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	hashData, ok := h.requests[requestId]
	if !ok {
		return errors.New("request not found")
	}

	if len(data) > 0 {
		hashData.Data = append(hashData.Data, data...)
	}

	hashData.CompletedWorkers++

	if hashData.CompletedWorkers == h.workersCount {
		hashData.Status = ReadyStatus
	}

	h.requests[requestId] = hashData

	return nil
}

func (h *HashUseCase) splitTask(requestId string, hash string, maxLength int) []WorkerTask {
	tasks := make([]WorkerTask, h.workersCount)

	totalCombinations := 0
	for i := 1; i <= maxLength; i++ {
		totalCombinations += int(math.Pow(float64(len(h.alphabet)), float64(i)))
	}

	for i := 0; i < h.workersCount; i++ {
		tasks[i] = WorkerTask{
			Hash:       hash,
			MaxLength:  maxLength,
			Alphabet:   h.alphabet,
			RequestId:  requestId,
			PartNumber: i,
			PartCount:  h.workersCount,
		}
	}

	return tasks
}

func (h *HashUseCase) startRequestTimeout(requestId string) {
	time.Sleep(time.Duration(h.ttl) * time.Second)

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if hashData, ok := h.requests[requestId]; ok && hashData.Status == InProgressStatus {
		if hashData.CompletedWorkers > 0 {
			hashData.Status = PartialReadyStatus
		} else {
			hashData.Status = ErrorStatus
		}
		h.requests[requestId] = hashData
		h.logger.Info("request timed out", slog.String("requestId", requestId))
	}
}
