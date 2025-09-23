package storage

import (
	"context"
	"encoding/json"
	stateUpdater "ksp-parser/stateUpdater/structures"
	"os"
)

type Storage interface {
	GetReminders(ctx context.Context) ([]Reminder, error)
	AddReminder(ctx context.Context, Task Reminder) error
	AddUpdateToReminder(ctx context.Context, ArticeId int, update stateUpdater.StateUpdate) error
}

type FileStorage struct {
	Path string
}

func (s *FileStorage) GetReminders(ctx context.Context) ([]Reminder, error) {
	data, err := os.ReadFile(s.Path)
	if err != nil {
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
