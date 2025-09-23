package stateupdater

import (
	"encoding/json"
	"fmt"
	"io"
	strcutures "ksp-parser/stateUpdater/structures"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type KspApi struct {
	client *http.Client
}

// NewKspApi создаёт клиента с cookie jar и таймаутом.
func NewKspApi() *KspApi {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: 5 * time.Second,
	}
	return &KspApi{client: client}
}

const (
	ItemEndpoint                 = "https://ksp.co.il/m_action/api/item/"
	BranchesAvailabilityEndpoint = "https://ksp.co.il/m_action/api/mlay/"
)

// GetAvailableBranches получает список филиалов, где есть товар в наличии.
func (a *KspApi) GetAvailableBranches(article int) ([]strcutures.Branch, error) {
	url := ItemEndpoint + fmt.Sprint(article)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; KspBot/1.0)")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res strcutures.ItemResultMain
	if err := json.Unmarshal(respBody, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal item response: %w", err)
	}

	branches, err := a.getBranchesInfo(res.Result.Data.Uinsql)
	if err != nil {
		return nil, err
	}

	var resBranches []strcutures.Branch
	for _, b := range branches {
		if b.Quantity > 0 {
			resBranches = append(resBranches, b)
		}
	}
	return resBranches, nil
}

// getBranchesInfo получает информацию о филиалах по uinsql.
func (a *KspApi) getBranchesInfo(uinsql string) ([]strcutures.Branch, error) {
	url := BranchesAvailabilityEndpoint + uinsql

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; KspBot/1.0)")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var branchesRes strcutures.BranchesInfoMain
	if err := json.Unmarshal(respBody, &branchesRes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal branches response: %w", err)
	}

	return strcutures.ConvertBranches(branchesRes.Result.Branches), nil
}
