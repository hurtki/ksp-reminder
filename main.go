package main

import (
	"ksp-parser/handlers"
	"ksp-parser/stateUpdater"
	"ksp-parser/storage"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// логгер

// зависимости
// Хранилище данных

// сущности
// Автоматический обновлятель состояния напоминаний
// Сам хендлер/api

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("KSP Reminder started")
	storage := storage.FileStorage{
		Path: "test.json",
	}
	stateUpdater := stateupdater.NewStateUpdater(&storage, logger)
	stateUpdater.UpdateInterval = time.Second * 10
	stateUpdater.Run()

	handler := handlers.NewTaskApiHandler(logger, &storage)

	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)

}
