package stateUpdater_kspApi

import (
	stateUpdaterStructures "ksp-parser/stateUpdater/structures"
	"strconv"
)

// Branches availability endpoint ksp
// structures represent what ksp api returns

type BranchesInfoMain struct {
	Result BranchesInfoResult `json:"result"`
}

type BranchesInfoResult struct {
	Branches map[string]KspApiBranch `json:"stores"`
	Shipment int                     `json:"shipment"`
}

type KspApiBranch struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Qnt  int    `json:"qnt"`
}

func (r *BranchesInfoResult) ToBranches() []stateUpdaterStructures.Branch {
	branches := make([]stateUpdaterStructures.Branch, 0, len(r.Branches))

	for key, b := range r.Branches {
		id, _ := strconv.Atoi(b.ID)
		branches = append(branches, stateUpdaterStructures.Branch{
			Id:       id,
			Name:     key,
			Quantity: b.Qnt,
		})
	}

	return branches
}
