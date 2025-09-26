package stateupdater

import (
	"context"
	"fmt"
	kspApi "ksp-parser/stateUpdater/kspApi"
	stateUpdaterStructures "ksp-parser/stateUpdater/structures"
	storage "ksp-parser/storage"
	"log/slog"
	"time"
)

type StateUpdater struct {
	UpdateInterval time.Duration
	storage.Storage
	kspApi *kspApi.KspApi
	logger *slog.Logger
}

func NewStateUpdater(storage storage.Storage, logger slog.Logger, kspApi *kspApi.KspApi) StateUpdater {
	return StateUpdater{
		Storage: storage,
		kspApi:  kspApi,
		logger:  logger.With("service", "StateUpdater"),
	}
}

func (u *StateUpdater) Run() {
	u.logger.Info("Started")
	go func() {
		for {
			u.logger.Info("Starting updating reminders")
			u.startUpdating()
			time.Sleep(u.UpdateInterval)
		}
	}()
}
// (u *StateUpdater) startUpdating() is getting reminders from storage and starting up workers, that update those reminders
func (u *StateUpdater) startUpdating() {
		reminders, err := u.Storage.GetReminders(context.TODO())
		if err != nil {
			u.logger.Error("cannot get reminders from storage", "error", err.Error())
			return
		}
		u.logger.Info("starting gorutines for updating reminders", "count", len(reminders))
		for i := range reminders {
			go u.updateReminder(reminders[i])
		}
}

// (u *StateUpdater)updateReminder() is a method function for (u *StateUpdater) startUpdating()
// (u *StateUpdater) updateReminder() is trying to access ksp api, generates stateUpdate and sends it to storage
func (u *StateUpdater) updateReminder(reminder storage.Reminder) {
	// starting filling up update
	stateUpdate := stateUpdaterStructures.StateUpdate{
		Time: time.Now(),
	}

	branches, err := u.kspApi.GetAvailableBranches(reminder.Article)

	// if no ability to get available branches from kspApi
	// we will send an update that there was an error getting them 
	if err != nil {
		u.logger.Error("cannot get branches via ksp api", "error", err.Error())
		stateUpdate.Error = err.Error()
		err = u.Storage.UpdateReminder(context.TODO(), reminder.Article, stateUpdate)
		if err != nil {
			u.logger.Error("cannot write update to reminder", "error", err.Error())
		}
		return
	}


	// trying to find the branch, that reminder is looking for
	for branchI := 0; branchI < len(branches); branchI++ {
		for reminderBranchI := 0; reminderBranchI < len(reminder.Branches); reminderBranchI++ {
			if branches[branchI].Id == reminder.Branches[reminderBranchI].Id {
				u.logger.Info("found branch on reminder article: " + fmt.Sprint(reminder.Article))
				stateUpdate.IsFound = true
				stateUpdate.BranchFoundOn = branches[branchI]
			}
		}
	}
	// logging if we didn't find anythink
	if !stateUpdate.IsFound {
		u.logger.Info("didn't find any branch on reminder article: " + fmt.Sprint(reminder.Article))
	}
	// sending update to the storage
	err = u.Storage.UpdateReminder(context.TODO(), reminder.Article, stateUpdate)
	if err != nil {
		u.logger.Error("cannot write update to reminder", "error", err.Error())
	}

}
