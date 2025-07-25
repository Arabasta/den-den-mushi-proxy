package puppet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

func (p *Client) getPuppetTaskUrl(t puppetTask) string {
	switch t {
	case TaskInjectPublicKey:
		return p.cfg.Endpoint + string(TaskInjectPublicKey)
	case TaskRemovePublicKey:
		return p.cfg.Endpoint + string(TaskRemovePublicKey)
	default:
		return ""
	}
}

func (p *Client) createPuppetRequest(t puppetTask, payload interface{}) (*http.Request, error) {
	url := p.getPuppetTaskUrl(t)
	if url == "" {
		return nil, fmt.Errorf("invalid puppet task: %s", t)
	}
	p.log.Debug("Creating puppet request", zap.String("url", url))

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// set headers todo: check if header ok
	req.Header.Set("Authorization", "token "+p.cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	var payloadMap map[string]interface{}
	_ = json.Unmarshal(body, &payloadMap)
	p.log.Debug("Puppet request created",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.Any("headers", req.Header),
		zap.Any("body", payloadMap))
	return req, nil
}

func (p *Client) httpPostAndResponse(r *http.Request) (*http.Response, error) {
	// call
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed to call endpoint: %w", err)
	}
	p.log.Debug("Puppet request successful", zap.Any("response", res))

	if res.StatusCode != http.StatusOK {
		p.log.Error("Puppet request failed", zap.Int("status_code", res.StatusCode))
		return nil, fmt.Errorf("status code not 200")
	}

	//// get response
	//respBody, err := io.ReadAll(res.Body)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to read response: %w", err)
	//}

	return res, nil
}
