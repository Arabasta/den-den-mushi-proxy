package puppet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type PuppetTask string

type taskBody struct {
	Environment string           `json:"environment"`
	Task        string           `json:"task"`
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

func (p *Client) createPuppetRequest(t PuppetTask, payload interface{}) (*http.Request, error) {
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

	var payloadMap map[string]interface{}
	_ = json.Unmarshal(body, &payloadMap)
	p.log.Debug("Puppet request created",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.Any("headers", req.Header),
		zap.Any("body", payloadMap))
	return req, nil
}

func (p *Client) callPuppetTask(task PuppetTask, payload interface{}) (*http.Response, error) {
	req, err := p.createPuppetRequest(task, payload)
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

var result struct {
	Job struct {
		Id     string `json:"id"`
		TaskID string `json:"name"`
	} `json:"job"`
}

func (p *Client) getPuppetTaskName(res *http.Response) string {
	// extract the taskname from the response
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			p.log.Error("Failed to close response body", zap.Error(err))
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		p.log.Error("Failed to read response body", zap.Error(err))
		return ""
	}
	if body == nil {
		p.log.Error("Puppet task returned no response")
		return ""
	}

	if err := json.Unmarshal(body, &result); err != nil {
		p.log.Error("Failed to unmarshal response body", zap.Error(err))
		return ""
	}
	if result.Job.TaskID == "" {
		p.log.Error("Puppet task returned empty task ID", zap.String("response", string(body)))
		return ""
	}

	return result.Job.TaskID
}
