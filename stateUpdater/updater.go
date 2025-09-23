package stateupdater

import (
	"context"
	"fmt"
	structures "ksp-parser/stateUpdater/structures"
	"ksp-parser/storage"
	"log/slog"
	"time"
)

type StateUpdater struct {
	UpdateInterval time.Duration
	storage.Storage
	kspApi *KspApi
	logger *slog.Logger
}

func NewStateUpdater(storage storage.Storage, logger *slog.Logger) StateUpdater {
	return StateUpdater{
		Storage: storage,
		kspApi:  NewKspApi(),
		logger:  logger,
	}
}

func (u *StateUpdater) Run() {
	u.logger.Info("State Updater Started")
	go u.startUpdating()
}

func (u *StateUpdater) startUpdating() {
	for {
		time.Sleep(u.UpdateInterval)
		reminders, err := u.Storage.GetReminders(context.TODO())
		if err != nil {
			u.logger.Error("cannot get reminders in state updater", "error", err.Error())
			continue
		}
		for i := range reminders {
			go u.updateReminder(reminders[i])
		}

	}

}

func (u *StateUpdater) updateReminder(reminder storage.Reminder) {
	branches, err := u.kspApi.GetAvailableBranches(reminder.Article)
	stateUpdate := structures.StateUpdate{
		Time: time.Now(),
	}

	if err != nil {
		u.logger.Error("cannot get branches via ksp api", "error", err.Error())
		stateUpdate.Error = err.Error()
		err = u.Storage.AddUpdateToReminder(context.TODO(), reminder.Article, stateUpdate)
		if err != nil {
			u.logger.Error("cannot write update to reminder", "error", err.Error())
		}
		return
	}

	for branchI := 0; branchI < len(branches); branchI++ {
		for reminderBranchI := 0; reminderBranchI < len(reminder.Branches); reminderBranchI++ {
			if branches[branchI].Id == reminder.Branches[reminderBranchI].Id {
				stateUpdate.IsFound = true
				stateUpdate.BranchFoundOn = branches[branchI]
			}
		}
	}
	err = u.Storage.AddUpdateToReminder(context.TODO(), reminder.Article, stateUpdate)
	if err != nil {
		u.logger.Error("cannot write update to reminder", "error", err.Error())
	}
	u.logger.Info("didn't find any branch on reminder article: " + fmt.Sprint(reminder.Article))
}
