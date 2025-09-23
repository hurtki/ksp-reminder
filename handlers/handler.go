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

func NewTaskApiHandler(logger *slog.Logger, storage storage.Storage) TaskApiHandler {
	return TaskApiHandler{
		logger:  logger,
		storage: storage,
	}
}

func (h *TaskApiHandler) ServeGet(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	h.logger.Info("GET request")
	tasks, err := h.storage.GetReminders(ctx)
	if err != nil {
		resp.Write([]byte(err.Error()))
		h.logger.Error(err.Error())
		return
	}
	res, err := json.Marshal(tasks)
	if err != nil {
		h.logger.Error(err.Error())
		resp.Write([]byte(err.Error()))
		return
	}
	resp.Write(res)
}

func (h *TaskApiHandler) ServePost(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	h.logger.Info("POST request")
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Error(err.Error())
		resp.WriteHeader(400)
		resp.Write([]byte(err.Error()))
		return
	}

	task := &storage.Reminder{}
	if err := json.Unmarshal(body, task); err != nil {
		h.logger.Error(err.Error())
		resp.WriteHeader(400)
		resp.Write([]byte(err.Error()))
		return
	}

	if err := h.storage.AddReminder(ctx, *task); err != nil {
		h.logger.Error(err.Error())
		resp.WriteHeader(400)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(201)
}

func (h *TaskApiHandler) ServeMethodNotAllowed(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(405)
}

func (h TaskApiHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	switch req.Method {
	case http.MethodGet:
		h.ServeGet(ctx, resp, req)
	case http.MethodPost:
		h.ServePost(ctx, resp, req)
	default:
		h.ServeMethodNotAllowed(ctx, resp, req)
	}
}
