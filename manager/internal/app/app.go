package app

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	mainrouter "manager/internal/adapters/http"
	proberouter "manager/internal/adapters/http/probe"
	"manager/internal/domain"
	"manager/internal/repository"
	httpserver "manager/pkg/http"
	"manager/pkg/logger"
)

func Run() {
	cfg, err := NewConfig()
	if err != nil {
		log.Panicf("failed to read config: %s", err)
	}

	loggerInstance := logger.NewLogger(cfg.LogLevel)

	repositoryInstance := repository.NewRepository(loggerInstance)

	hashUseCase := domain.NewHashUseCase(
		cfg.WorkerCount,
		strings.Split(cfg.WorkerURLs, ","),
		cfg.Alphabet,
		cfg.TTL,
		repositoryInstance,
		loggerInstance,
	)

	mainRouter := mainrouter.NewRouter(hashUseCase, loggerInstance)
	mainServer := httpserver.New(cfg.MainServerPort, mainRouter)

	probeRouter := proberouter.NewRouter(loggerInstance)
	probeServer := httpserver.New(cfg.ProbeServerPort, probeRouter)

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		loggerInstance.Info("starting main server")
		if err = mainServer.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			loggerInstance.Error("main server error: %s", err)
		}
	}()

	go func() {
		defer wg.Done()
		loggerInstance.Info("starting probe server")
		if err = probeServer.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			loggerInstance.Error("probe server error: %s", err)
		}
	}()

	<-stopCh
	loggerInstance.Info("shutting down gracefully")

	if err = mainServer.Shutdown(); err != nil {
		loggerInstance.Error("main server shutdown error: %v", err)
	}

	if err = probeServer.Shutdown(); err != nil {
		loggerInstance.Error("probe server shutdown error: %v", err)
	}

	wg.Wait()

	loggerInstance.Info("servers stopped gracefully")
}
