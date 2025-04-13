package repository

import (
	"bytes"
	"encoding/xml"
	"errors"
	"log/slog"
	"net/http"
)

const managerSolveTaskEndpoint = "/internal/api/manager/hash/crack/request"

type Repository struct {
	managerURL string
	logger     *slog.Logger
}

func NewRepository(managerURL string, logger *slog.Logger) *Repository {
	return &Repository{
		managerURL: managerURL,
		logger:     logger,
	}
}

func (r *Repository) SendCrackedTaskToManager(requestId string, words []string) error {
	workerResponse := WorkerResponse{
		RequestId: requestId,
		Words:     words,
	}

	xmlData, err := xml.Marshal(workerResponse)
	if err != nil {
		r.logger.Error("failed to marshal task to XML", "error", err)
		return err
	}

	url := r.managerURL + managerSolveTaskEndpoint

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(xmlData))
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

	r.logger.Info("solve task successfully sent to manager", "managerURL", r.managerURL)

	return nil
}
