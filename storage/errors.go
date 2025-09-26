package storage

import "errors"

var (
	ErrNotFoundInStorage = errors.New("not found in storage")
	ErrReminderAlreadyExists = errors.New("reminder you are trying to add already exists")
)
