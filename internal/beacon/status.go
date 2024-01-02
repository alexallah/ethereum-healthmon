package beacon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type syncInfoResponse struct {
	Data *syncInfo `json:"data"`
}

func buildUrl(addr string, endpoint string) string {
	return fmt.Sprintf("%s/%s", addr, endpoint)
}

func getSyncing(addr string, client *http.Client) (*syncInfo, error) {

	url := buildUrl(addr, "eth/v1/node/syncing")
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	syncInfo := new(syncInfoResponse)
	err = json.Unmarshal(body, syncInfo)
	if err != nil {
		return nil, err
	}
	return syncInfo.Data, nil
}

func isHealthy(addr string, client *http.Client) error {
	url := buildUrl(addr, "eth/v1/node/health")
	res, err := client.Get(url)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("health status %d", res.StatusCode)
	}
	return nil
}

func isReady(addr string, client *http.Client) error {
	// make sure we are not syncing
	syncInfo, err := getSyncing(addr, client)
	if err != nil {
		return err
	}

	if err := checkSyncInfo(syncInfo); err != nil {
		return err
	}

	// health call
	if err := isHealthy(addr, client); err != nil {
		return err
	}

	return nil
}
