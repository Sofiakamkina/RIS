package domain

import (
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type HashUseCase struct {
	mutex          sync.Mutex
	requests       map[string]HashData
	workersCount   int
	workersURL     []string
	alphabet       string
	ttl            int
	httpRepository IHTTPRepository
	sqlRepository  ISQLRepository
	logger         *slog.Logger
}

func NewHashUseCase(
	workersCount int,
	workersURL []string,
	alphabet string,
	ttl int,
	httpRepository IHTTPRepository,
	sqlRepository ISQLRepository,
	logger *slog.Logger,
) *HashUseCase {
	return &HashUseCase{
		requests:       make(map[string]HashData),
		workersCount:   workersCount,
		workersURL:     workersURL,
		alphabet:       alphabet,
		ttl:            ttl,
		httpRepository: httpRepository,
		sqlRepository:  sqlRepository,
		logger:         logger,
	}
}

func (h *HashUseCase) Init() error {
	crackHashes, err := h.sqlRepository.GetAll()
	if err != nil {
		return err
	}

	for _, crackHash := range crackHashes {
		requestId := &crackHash.RequestId
		_, err = h.CrackHash(requestId, crackHash.Hash, crackHash.MaxLength)
		if err != nil {
			return err
		}
		h.logger.Info("init hash successfully", "request_id", requestId, "hash", requestId)
	}

	return nil
}

func (h *HashUseCase) CrackHash(requestId *string, hash string, maxLength int) (HashData, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if requestId == nil {
		id := uuid.New().String()
		requestId = &id
	}

	hashData := HashData{
		RequestId: *requestId,
		Status:    InProgressStatus,
	}

	h.requests[*requestId] = hashData

	if err := h.sqlRepository.Save(*requestId, hash, maxLength); err != nil {
		return hashData, err
	}
	h.logger.Info("hash save successfully", "request_id", *requestId, "hash", hash, "max_length", maxLength)

	tasks := h.splitTask(*requestId, hash, maxLength)

	successChan := make(chan bool, h.workersCount)

	for i, task := range tasks {
		go func(i int, task WorkerTask) {
			if err := h.httpRepository.SendTaskToWorker(h.workersURL[i], task); err != nil {
				h.logger.Error("failed to send task to worker", slog.String("workerURL", h.workersURL[i]), slog.String("error", err.Error()))
				for j := 0; j < h.workersCount; j++ {
					if j == i {
						continue
					}
					if err = h.httpRepository.SendTaskToWorker(h.workersURL[j], task); err == nil {
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

	go h.startRequestTimeout(*requestId)

	return hashData, nil
}

func (h *HashUseCase) GetStatus(requestId string) (HashData, int, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	response, ok := h.requests[requestId]
	if !ok {
		return HashData{}, 0, errors.New("request not found")
	}

	return response, h.workersCount, nil
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
		h.logger.Info("hash cracked", slog.String("requestId", requestId), slog.String("data", strings.Join(hashData.Data, ",")))
		if err := h.sqlRepository.Remove(requestId); err != nil {
			h.logger.Error("failed to remove request", slog.String("requestId", requestId))
		}
	} else {
		hashData.Status = PartialReadyStatus
	}

	h.requests[requestId] = hashData

	return nil
}

func (h *HashUseCase) splitTask(requestId string, hash string, maxLength int) []WorkerTask {
	tasks := make([]WorkerTask, h.workersCount)
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
		if hashData.CompletedWorkers == 0 {
			hashData.Status = ErrorStatus
		}
		h.requests[requestId] = hashData
		h.logger.Info("request timed out", slog.String("requestId", requestId), slog.String("status", string(hashData.Status)))
	}
}
