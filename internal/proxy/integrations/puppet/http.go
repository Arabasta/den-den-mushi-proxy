package puppet

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

func (p *Client) httpPostAndResponse(r *http.Request) (*http.Response, error) {
	// call
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed to call endpoint: %w", err)
	}
	p.log.Debug("Http request successful", zap.Any("response", res))

	//if true {
	//	p.log.Error("Http request failed", zap.Int("status_code", res.StatusCode))
	//	return nil, fmt.Errorf("status code not 200")
	//}

	//// get response
	//respBody, err := io.ReadAll(res.Body)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to read response: %w", err)
	//}

	return res, nil
}
