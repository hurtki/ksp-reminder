package stateUpdater_kspApi

import (
	"encoding/json"
	"fmt"
	"io"
	structures "ksp-parser/stateUpdater/structures"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type KspApi struct {
	client *http.Client
}

func NewKspApi(timeout time.Duration) *KspApi {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: timeout,
	}
	return &KspApi{client: client}
}

const (
	ItemEndpoint                 = "https://ksp.co.il/m_action/api/item/"
	BranchesAvailabilityEndpoint = "https://ksp.co.il/m_action/api/mlay/"
)

func (a *KspApi) RequestApi(method string, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, KspApiError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

// GetAvailableBranches retrieves available branches for exact product
func (a *KspApi) GetAvailableBranches(article int) ([]structures.Branch, error) {
	url := ItemEndpoint + fmt.Sprint(article)

	respBody, err := a.RequestApi("GET", url, nil)

	if err != nil {
		return nil, err
	}

	var res ItemResultMain
	if err := json.Unmarshal(respBody, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal item response: %w", err)
	}
	
	branches, err := a.getBranchesInfo(res.Result.Data.Uinsql)
	if err != nil {
		return nil, err
	}

	var resBranches []structures.Branch
	for _, b := range branches {
		if b.Qnt > 0 {
			resBranches = append(resBranches, b.ToInternalBranch())
		}
	}
	return resBranches, nil
}

// getBranchesInfo получает информацию о филиалах по uinsql.
func (a *KspApi) getBranchesInfo(uinsql string) ([]KspApiBranch, error) {
	url := BranchesAvailabilityEndpoint + uinsql

	respBody, err := a.RequestApi("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var branchesRes BranchesInfoMain
	if err := json.Unmarshal(respBody, &branchesRes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal branches response: %w", err)
	}

	return branchesRes.Result.ToBranches(), nil
}
