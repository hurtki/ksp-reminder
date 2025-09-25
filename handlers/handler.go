package handlers

import (
	"context"
	"encoding/json"
	"io"
	"ksp-parser/storage"
	"log/slog"
	"net/http"
)

type TaskApiHandler struct {
	logger  *slog.Logger
	storage storage.Storage
}

func NewTaskApiHandler(logger slog.Logger, storage storage.Storage) TaskApiHandler {
	return TaskApiHandler{
		logger:  logger.With("service", "HTTP-Handler"),
		storage: storage,
	}
}

func (h *TaskApiHandler) ServeGet(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	reminders, err := h.storage.GetReminders(ctx)
	if err != nil {
		resp.Write([]byte(err.Error()))
		h.logger.Error("failed to get reminders", "error", err)
		return
	}
	res, err := json.Marshal(reminders)
	if err != nil {
		h.logger.Error("failed to marshal reminders", "error", err, "reminders", reminders)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.Write(res)
}

func (h *TaskApiHandler) ServePost(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("failed to read request body", "error", err)
		resp.WriteHeader(400)
		resp.Write([]byte(err.Error()))
		return
	}

	reminder := &storage.Reminder{}
	if err := json.Unmarshal(body, reminder); err != nil {
		h.logger.Error(err.Error())
		resp.WriteHeader(400)
		resp.Write([]byte(err.Error()))
		return
	}

	if err := h.storage.AddReminder(ctx, *reminder); err != nil {
		h.logger.Error("failed to add reminder", "error", err, "reminder", reminder)
		resp.WriteHeader(400)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(201)
}

func (h *TaskApiHandler) ServeMethodNotAllowed(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	h.logger.Error("Incoming request Method NOT ALLOWED",
		"method", req.Method,
		"path", req.URL.Path,
		"remote", req.RemoteAddr,
	)
	resp.WriteHeader(405)
}

func (h TaskApiHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	// logging request info

	h.logger.Info("incoming request",
		"method", req.Method,
		"path", req.URL.Path,
		"remote", req.RemoteAddr,
	)

	switch req.Method {
	case http.MethodGet:
		h.ServeGet(ctx, resp, req)
	case http.MethodPost:
		h.ServePost(ctx, resp, req)
	default:
		h.ServeMethodNotAllowed(ctx, resp, req)
	}
}
