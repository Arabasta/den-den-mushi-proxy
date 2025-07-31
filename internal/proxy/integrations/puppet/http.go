package puppet

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

func (p *Client) httpPostAndResponse(r *http.Request) (*http.Response, error) {
	p.log.Debug("Sending Puppet request", zap.String("URL", r.URL.String()))
	p.log.Debug("Request Method", zap.String("Method", r.Method))
	p.log.Debug("Request Headers", zap.Any("Headers", r.Header))
	p.log.Debug("Request Body", zap.Any("Body", r.Body))

	// call
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed to call endpoint: %w", err)
	}
	//p.log.Debug("Http request successful", zap.Any("response", res))

	//// get response
	//respBody, err := io.ReadAll(res.Body)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to read response: %w", err)
	//}

	return res, nil
}
