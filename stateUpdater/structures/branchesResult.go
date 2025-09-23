package stateUpdaterStructures

type KspApiBranch struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Qnt  int    `json:"qnt"`
}

type BranchesInfoResult struct {
	Branches map[string]KspApiBranch `json:"stores"`
	Shipment int                     `json:"shipment"`
}

type BranchesInfoMain struct {
	Result BranchesInfoResult `json:"result"`
}
