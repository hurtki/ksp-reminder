package storage

import (
	"context"
	"encoding/json"
	stateUpdater "ksp-parser/stateUpdater/structures"
	"log/slog"
	"os"
	"sync"
)

type Storage interface {
	GetReminders(ctx context.Context) ([]Reminder, error)
	AddReminderIfNotExists(ctx context.Context, Task Reminder) error
	UpdateReminder(ctx context.Context, ArticeId int, update stateUpdater.StateUpdate) error
}

type FileStorage struct {
	Path   string
	logger *slog.Logger
	mutex sync.Mutex
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
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.readReminders()
}

func (s *FileStorage) AddReminderIfNotExists(ctx context.Context, ReminderToAdd Reminder) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	reminders, err := s.readReminders()
	if err != nil {
		return err
	}
	// checking if Article already exists in storage
	for _, r := range reminders {
		if r.Article == ReminderToAdd.Article {
			return ErrReminderAlreadyExists
		}
	}
	reminders = append(reminders, ReminderToAdd)
	
	
	
	return s.writeReminders(reminders)
}

func (s *FileStorage) UpdateReminder(ctx context.Context, articleID int, update stateUpdater.StateUpdate) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	reminders, err := s.readReminders()

	if err != nil {
		return err
	}

	for i := range reminders {
		if reminders[i].Article != articleID {
			continue
		}

		updates := reminders[i].Updates
		if len(updates) > 0 {
			last := &updates[len(updates)-1]

			// if status of update is the same — so updateing only time and BranchFoundOn
			if last.Error == update.Error && last.IsFound == update.IsFound {
				last.Time = update.Time
				last.BranchFoundOn = update.BranchFoundOn
			} else {
				updates = append(updates, update)
			}
		} else {
			// if there is no updates just add a new one
			updates = append(updates, update)
		}
		
		reminders[i].Updates = updates

		return s.writeReminders(reminders)
	}

	return ErrNotFoundInStorage
}

// readReminders() is a raw method to read reminders from file 
// You need to use s.mutex in functions that use this
func (s *FileStorage) readReminders() ([]Reminder, error) {
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
// readReminders() is a raw method to write reminders to file 
// You need to use s.mutex in functions that use this
func (s *FileStorage) writeReminders(reminders []Reminder) error {
	data, err := json.Marshal(reminders)
	if err != nil {
		return err
	}
	return os.WriteFile(s.Path, data, 0644)
}