package puppet

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type puppetTask string

const (
	TaskInjectPublicKey puppetTask = "inject_public_key"
	TaskRemovePublicKey puppetTask = "remove_public_key"
)

func (p *Client) callPuppetTask(task puppetTask, payload interface{}) (*http.Response, error) {
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
