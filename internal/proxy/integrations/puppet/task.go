package puppet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type puppetTask string

type taskBody struct {
	Environment string           `json:"environment"`
	Task        puppetTask       `json:"task"`
	Params      sshKeyTaskParams `json:"params"`
	Scope       taskScope        `json:"scope"`
}

type sshKeyTaskParams struct {
	PublicKey    string `json:"public_key"`
	OSUser       string `json:"os_user"`
	ConnectionID string `json:"connection_id"`
}

type taskScope struct {
	Nodes []string `json:"nodes"`
}

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
	//req.Header.Set("Authorization", "token "+p.cfg.Token)
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

func (p *Client) callPuppetTask(task puppetTask, payload taskBody) (*http.Response, error) {
	req, err := p.createPuppetRequest(task, payload)
	if err != nil {
		return nil, err
	}

	var resp *http.Response
	for i := 1; i <= p.cfg.RetryAttempts; i++ {
		resp, err = p.httpPostAndResponse(req)
		if err != nil {
			p.log.Error(fmt.Sprintf("Failed to send Puppet request. Attempt %d of %d", i, p.cfg.RetryAttempts), zap.Error(err))
			if i == p.cfg.RetryAttempts {
				return nil, err
			}
			time.Sleep(p.cfg.TaskRetrySeconds * time.Second)
		}
	}

	p.log.Debug("Puppet task completed", zap.String("task", string(task)), zap.Any("response", resp))
	return resp, nil
}
