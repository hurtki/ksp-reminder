package stateUpdater_kspApi

import "fmt"

type KspApiError struct {
	StatusCode int
	Status     string
}

func (e KspApiError) Error() string {
	return fmt.Sprintf("API error: %d %s", e.StatusCode, e.Status)
}
