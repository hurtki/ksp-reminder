package storage

import (
	stateUpdater "ksp-parser/stateUpdater/structures"
)

type Reminder struct {
	Article  int
	Branches []stateUpdater.Branch
	Updates  []stateUpdater.StateUpdate
}
