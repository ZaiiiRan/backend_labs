package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	config "github.com/ZaiiiRan/backend_labs/order-service/internal/config/settings"
	"github.com/ZaiiiRan/backend_labs/order-service/pkg/api/dto/v1"
)

type OmsHttpClient struct {
	baseAddress string
	httpClient  *http.Client
}

func NewOmsHttpClient(cfg config.HttpClientSettings) *OmsHttpClient {
	return &OmsHttpClient{
		baseAddress: cfg.BaseAdress,
		httpClient:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *OmsHttpClient) LogOrder(req *dto.V1CreateAuditLogOrderRequest) (*dto.V1CreateAuditLogOrderResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := c.baseAddress + "/api/v1/audit-log/order/batch-create"

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("post request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var logResp dto.V1CreateAuditLogOrderResponse
	if err := json.Unmarshal(bodyBytes, &logResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &logResp, nil
}
