package clients

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-logr/logr"
)

const (
	httpBasicAuthKey     = "Basic"
	errCreatingRequest   = "Failed creating request of method"
	httpAuthorizationKey = "Authorization"
	errDoingRequest      = "Failed making the request"
	errClosingBody       = "Failed to close body"
	errReadingRequest    = "Error in reading the request"
)

// SendRequest creates a http request and returns the response body.
func SendRequest(url string, authHeader string, logger logr.Logger, httpClient *http.Client) ([]byte, error) {
	httpBasicAuthPrefix := fmt.Sprintf("%s %s", httpBasicAuthKey, authHeader)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		logger.Error(err, "%s", errCreatingRequest, http.MethodPost, err)
		return nil, fmt.Errorf("%s %q: %w", errCreatingRequest, http.MethodPost, err)
	}
	req.Header.Set(httpAuthorizationKey, httpBasicAuthPrefix)

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error(err, errDoingRequest)
		return nil, fmt.Errorf("%s", errDoingRequest)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err, errClosingBody)
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err, errReadingRequest)
		return nil, fmt.Errorf("%s", errReadingRequest)
	}
	return body, nil
}
