package stateUpdaterStructures

import "time"

type StateUpdate struct {
	Time          time.Time
	IsFound       bool
	BranchFoundOn Branch
	Error         string
}
