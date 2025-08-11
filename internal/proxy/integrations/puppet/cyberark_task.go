package puppet

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

// todo: refactor, rush for demo as usual

const TaskCyberarkDrawKey PuppetTask = "cyberark_draw_key"

type taskBody2 struct {
	Environment string                    `json:"environment"`
	Task        string                    `json:"task"`
	Params      cyberarkdrawkeyTaskParams `json:"params"`
	Scope       taskScope                 `json:"scope"`
}

type cyberarkdrawkeyTaskParams struct {
	CybidA  string `json:"cybid_a"`
	CybidB  string `json:"cybid_b"`
	SafeA   string `json:"safe_a"`
	SafeB   string `json:"safe_b"`
	ObjectA string `json:"object_a"`
	ObjectB string `json:"object_b"`
}

func (p *Client) DrawCyberarkKey(object string, serverFQDN string) (string, error) {
	params := taskBody2{
		Environment: p.cfg.TaskEnvironment,
		Task:        p.Cfgtmp.PuppetTasks.CyberarkPasswordDraw.TaskName,
		Params: cyberarkdrawkeyTaskParams{
			CybidA:  p.Cfgtmp.PuppetTasks.CyberarkPasswordDraw.CybidA,
			CybidB:  p.Cfgtmp.PuppetTasks.CyberarkPasswordDraw.CybidB,
			SafeA:   p.Cfgtmp.PuppetTasks.CyberarkPasswordDraw.SafeA,
			SafeB:   p.Cfgtmp.PuppetTasks.CyberarkPasswordDraw.SafeB,
			ObjectA: object,
			ObjectB: object,
		},
		Scope: taskScope{
			Nodes: []string{p.cfg.TaskNode},
		},
	}

	p.log.Debug("Calling puppet task to draw key with params", zap.Any("params", params))

	res, err := p.callPuppetTask(TaskCyberarkDrawKey, params)
	if err != nil {
		return "", fmt.Errorf("failed to call puppet task to draw cyberark key: %w", err)
	}

	taskName := p.getPuppetTaskName(res)
	p.log.Debug("Cyberark draw key task called", zap.String("Task name", taskName))
	if taskName == "" {
		return "", fmt.Errorf("failed to get puppet task name from response, task may not have been created")
	}

	password, err := p.queryOrchestratorJobForCyberarkPasswordWhatever(taskName)
	if err != nil {
		return "", fmt.Errorf("failed to query orchestrator job: %w", err)
	}

	return password, nil
}

func (p *Client) queryOrchestratorJobForCyberarkPasswordWhatever(jobId string) (string, error) {
	p.log.Debug("Query orchestrator job for cyberark password whatever started, sleeping before query")
	time.Sleep(p.Cfgtmp.PuppetTasks.QueryJobs.WaitBeforeQuerySeconds * time.Second)

	url := fmt.Sprintf("%s/%s/nodes", p.Cfgtmp.PuppetTasks.QueryJobs.OrchestratorEndpoint, jobId)
	p.log.Debug("Querying orchestrator job for cyberark password", zap.String("url", url))

	type response struct {
		Items []struct {
			Result struct {
				Status string `json:"status"`
				Stdout string `json:"stdout"`
			} `json:"result"`
		} `json:"items"`
	}

	var jobResp response

	for i := 0; i < p.Cfgtmp.PuppetTasks.QueryJobs.MaxQueryAttempts; i++ {
		p.log.Debug("Querying orchestrator job for cyberark password", zap.String("jobId", jobId),
			zap.Int("attempt", i+1), zap.Int("maxAttempts", p.Cfgtmp.PuppetTasks.QueryJobs.MaxQueryAttempts))
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return "", fmt.Errorf("failed to build request: %w", err)
		}

		resp, err := p.httpPostAndResponse(req)
		if err != nil {
			p.log.Warn("orchestrator query failed", zap.Int("attempt", i+1), zap.Error(err))
			time.Sleep(p.Cfgtmp.PuppetTasks.QueryJobs.QueryIntervalSeconds * time.Second)
			continue
		}
		if resp == nil {
			p.log.Warn("orchestrator returned nil response", zap.Int("attempt", i+1))
			time.Sleep(p.Cfgtmp.PuppetTasks.QueryJobs.QueryIntervalSeconds * time.Second)
			continue
		}

		p.log.Debug("Query orchestrator job response status", zap.Int("status", resp.StatusCode))

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			p.log.Warn("failed to read orchestrator response body", zap.Error(err))
			time.Sleep(p.Cfgtmp.PuppetTasks.QueryJobs.QueryIntervalSeconds * time.Second)
			continue
		}

		if err := json.Unmarshal(body, &jobResp); err != nil {
			p.log.Error("Failed to unmarshal orchestrator response", zap.String("response", string(body)), zap.Error(err))
			time.Sleep(p.Cfgtmp.PuppetTasks.QueryJobs.QueryIntervalSeconds * time.Second)
			continue
		}

		// check status of the job
		if len(jobResp.Items) > 0 {
			result := jobResp.Items[0].Result
			if result.Status == "success" && result.Stdout != "" {
				p.log.Debug("Job succeeded, returning password", zap.String("password", result.Stdout))
				// return the password
				return result.Stdout, nil
			}
			// if status not "success", keep polling. who designed this?
		}

		time.Sleep(p.Cfgtmp.PuppetTasks.QueryJobs.QueryIntervalSeconds * time.Second)
	}

	return "", fmt.Errorf("job %s did not succeed or had no stdout after %d attempts",
		jobId, p.Cfgtmp.PuppetTasks.QueryJobs.MaxQueryAttempts)
}
