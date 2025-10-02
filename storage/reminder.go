package storage

import (
	stateUpdater "ksp-parser/stateUpdater/structures"
)

type Reminder struct {
	Article     int
	BranchesIDs []int
	Updates     []stateUpdater.StateUpdate
}
