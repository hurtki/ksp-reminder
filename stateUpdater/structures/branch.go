package stateUpdaterStructures

import "strconv"

type Branch struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"qnt"`
}

func ConvertBranches(apiBranches map[string]KspApiBranch) []Branch {
	branches := make([]Branch, 0, len(apiBranches))

	for key, b := range apiBranches {
		id, _ := strconv.Atoi(b.ID)
		branches = append(branches, Branch{
			Id:       id,
			Name:     key,
			Quantity: b.Qnt,
		})
	}

	return branches
}
