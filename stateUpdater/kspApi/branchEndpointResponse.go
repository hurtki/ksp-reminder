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
	Id   string `json:"id"`
	Name string `json:"name"`
	Qnt  int    `json:"qnt"`
}

func (b *KspApiBranch) ToInternalBranch() stateUpdaterStructures.Branch {
	id, _ := strconv.Atoi(b.Id)
	return stateUpdaterStructures.Branch{
		Id:   id,
		Name: b.Name,
	}
}

func (r *BranchesInfoResult) ToBranches() []KspApiBranch {
	branches := make([]KspApiBranch, 0, len(r.Branches))

	for key, b := range r.Branches {
		branches = append(branches, KspApiBranch{
			Id:   b.Id,
			Name: key,
		})
	}

	return branches
}
