package puppet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

func (p *Client) createPuppetRequest2(t PuppetTask, payload taskBody2) (*http.Request, error) {
	url := p.cfg.Endpoint
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

	req.Header.Set("Content-Type", "application/json")

	//var payloadMap map[string]interface{}
	//_ = json.Unmarshal(body, &payloadMap)
	p.log.Debug("Puppet request created",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.Any("headers", req.Header),
		zap.ByteString("body", body))
	return req, nil
}

func (p *Client) callPuppetTask2(task PuppetTask, payload taskBody2) (*http.Response, error) {
	req, err := p.createPuppetRequest2(task, payload)
	if err != nil {
		return nil, err
	}

	resp, err := p.httpPostAndResponse(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("puppet returned nil response")
	}

	p.log.Debug("Puppet task completed", zap.String("task", string(task)))
	return resp, nil
}
