package repository

import (
	"bytes"
	"encoding/xml"
	"errors"
	"log/slog"
	"net/http"

	"manager/internal/domain"
)

const workerTaskEndpoint = "/internal/api/worker/hash/crack/task"

type Repository struct {
	logger *slog.Logger
}

func NewRepository(logger *slog.Logger) *Repository {
	return &Repository{
		logger: logger,
	}
}

func (r *Repository) SendTaskToWorker(workerURL string, task domain.WorkerTask) error {
	workerRequest := WorkerRequest{
		RequestId:  task.RequestId,
		Hash:       task.Hash,
		Alphabet:   task.Alphabet,
		MaxLength:  task.MaxLength,
		PartNumber: task.PartNumber,
		PartCount:  task.PartCount,
	}

	xmlData, err := xml.Marshal(workerRequest)
	if err != nil {
		r.logger.Error("failed to marshal task to XML", "error", err)
		return err
	}

	url := workerURL + workerTaskEndpoint

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(xmlData))
	if err != nil {
		r.logger.Error("failed to create HTTP request", "error", err)
		return err
	}

	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		r.logger.Error("failed to send HTTP request", "error", err)
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		r.logger.Error("received non-OK HTTP status", "status", resp.StatusCode)
		return errors.New("received non-OK HTTP status")
	}

	r.logger.Info("task successfully sent to worker", "workerURL", workerURL)

	return nil
}
