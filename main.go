package main

import (
	"ksp-parser/handlers"
	"ksp-parser/stateUpdater"
	"ksp-parser/storage"
	"log/slog"
	"net/http"
	"os"
	"time"

	kspApi "ksp-parser/stateUpdater/kspApi"
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
	storage := storage.NewFileStorage("data.json", *logger)
	if err := storage.Init(); err != nil {
		logger.Error("Storage wasn't initialized, leaving programm", "error", err)
	}
	stateUpdater := stateupdater.NewStateUpdater(&storage, *logger, kspApi.NewKspApi(5*time.Second))
	stateUpdater.UpdateInterval = time.Second * 10
	stateUpdater.Run()

	handler := handlers.NewTaskApiHandler(*logger, &storage)

	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)

}
