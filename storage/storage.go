package storage

import (
	"context"
	"encoding/json"
	stateUpdater "ksp-parser/stateUpdater/structures"
	"log/slog"
	"os"
)

type Storage interface {
	GetReminders(ctx context.Context) ([]Reminder, error)
	AddReminder(ctx context.Context, Task Reminder) error
	AddUpdateToReminder(ctx context.Context, ArticeId int, update stateUpdater.StateUpdate) error
}

type FileStorage struct {
	Path   string
	logger *slog.Logger
}

func NewFileStorage(path string, logger slog.Logger) FileStorage {
	return FileStorage{
		Path:   path,
		logger: logger.With("service", "FileStorage"),
	}
}

func (s *FileStorage) Init() error {
	s.logger.Info("Started Initialization of file storage", "path", s.Path)
	_, err := os.Stat(s.Path)
	if os.IsNotExist(err) {

		data, _ := json.Marshal([]Reminder{})
		err = os.WriteFile(s.Path, data, 0644)
		if err != nil {
			return err
		}
		s.logger.Info("Created new storage file", "path", s.Path)
		return nil
	}
	if err != nil {
		return err
	}
	s.logger.Info("File Storage already exists", "path", s.Path)
	return nil
}


func (s *FileStorage) GetReminders(ctx context.Context) ([]Reminder, error) {

	data, err := os.ReadFile(s.Path)
	if err != nil {
		s.logger.Warn("Filepath storage didn't find a storage on path: " + s.Path)
		data, _ = json.Marshal([]Reminder{})
	}
	reminders := []Reminder{}
	err = json.Unmarshal(data, &reminders)
	if err != nil {
		return nil, err
	}
	return reminders, nil
}

func (s *FileStorage) AddReminder(ctx context.Context, TaskToAdd Reminder) error {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		s.logger.Warn("Filepath storage didn't find a storage on path: " + s.Path)
		data, _ = json.Marshal([]Reminder{})
	}
	reminders := &[]Reminder{}
	_ = json.Unmarshal(data, reminders)
	*reminders = append(*reminders, TaskToAdd)
	updatedData, err := json.Marshal(reminders)
	if err != nil {
		return err
	}
	err = os.WriteFile(s.Path, updatedData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *FileStorage) AddUpdateToReminder(ctx context.Context, ArticeId int, update stateUpdater.StateUpdate) error {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		data, _ = json.Marshal([]Reminder{})
	}
	reminders := []Reminder{}

	_ = json.Unmarshal(data, &reminders)

	for i := 0; i < len(reminders); i++ {
		if reminders[i].Article == ArticeId {
			reminders[i].Updates = append(reminders[i].Updates, update)
			updatedData, err := json.Marshal(reminders)
			if err != nil {
				return err
			}
			err = os.WriteFile(s.Path, updatedData, 0644)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return ErrNotFoundInStorage
}
