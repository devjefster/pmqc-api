package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pmqc-api/internal/config"
	"pmqc-api/internal/models"
)

type FetchResult struct {
	Year  int
	Month int
	Data  interface{}
	Error error
}

func FetchPMQCData(year int, month int) (*models.PMQCData, error) {
	url := fmt.Sprintf("https://www.gov.br/anp/pt-br/centrais-de-conteudo/dados-abertos/arquivos/pmqc/%d/%d-%02d-pmqc.json", year, year, month)

	config.Logger.Infof("🌐 Fetching: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		config.Logger.Errorf("🚨 HTTP request failed: %v", err)
		return nil, fmt.Errorf("error fetching %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		config.Logger.Warnf("⚠️ Received status %d from %s", resp.StatusCode, url)
		return nil, fmt.Errorf("received status %d from %s", resp.StatusCode, url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		config.Logger.Errorf("🚨 Error reading response body: %v", err)
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var pmqcData models.PMQCData
	err = json.Unmarshal(body, &pmqcData)
	if err != nil {
		config.Logger.Errorf("🚨 JSON unmarshalling error: %v", err)
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	config.Logger.Infof("✅ Successfully fetched PMQC data for %d-%02d!", year, month)
	return &pmqcData, nil
}
