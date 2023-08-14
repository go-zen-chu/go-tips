package prometheus

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	retryCount = 3
)

type PrometheusClient interface {
	//TBD
}

type prometheusClient struct {
	client http.Client
}

func NewPrometheusClient(httpClient http.Client) PrometheusClient {
	return &prometheusClient{
		client: httpClient,
	}
}

type PrometheusResponse struct {
	Status string      `json:"status"`
	Data   MetricsData `json:"data"`
}

type MetricsData struct {
	ResultType string   `json:"resultType"`
	Results    []Result `json:"result"`
}

type Result struct {
	Metric map[string]string `json:"metric"`
	Value  []any             `json:"value"`
}

// RequestPrometheus は promURL に対して promQL のメトリクス取得を行う
func (p *prometheusClient) RequestPrometheus(promURL, promQL string) (*PrometheusResponse, error) {
	req, err := http.NewRequest("GET", promURL, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("query", promQL)
	req.URL.RawQuery = q.Encode()

	var resp *http.Response
	for i := 1; i <= retryCount; i++ {
		resp, err = p.client.Do(req)
		if err != nil {
			// 並列処理をすると、環境によってたまに名前解決に失敗するため、リトライさせる
			if ne, ok := err.(*net.OpError); ok && ne.Op == "dial" {
				log.Printf("failed to resolve %s, retrying: %d", promURL, i)
				time.Sleep(time.Duration(i) * time.Second)
				continue
			}
			return nil, err
		}
	}
	bd, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not access to prometheus: %w", err)
	}
	var pres PrometheusResponse
	if err := json.Unmarshal(bd, &pres); err != nil {
		return nil, fmt.Errorf("invalid json (%s): %w", string(bd), err)
	}
	if len(pres.Data.Results) == 0 {
		return &pres, errors.New("no result")
	}
	return &pres, nil
}
