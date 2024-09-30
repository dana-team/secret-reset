package clients

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-logr/logr"
)

const (
	httpBasicAuthKey     = "Basic"
	httpAuthorizationKey = "Authorization"
)

const (
	errCreatingRequest = "failed creating request of method"
	errDoingRequest    = "failed making request"
	errClosingBody     = "failed closing body"
	errReadingRequest  = "failed reading request"
)

// SendRequest creates a http request and returns the response body.
func SendRequest(url string, authHeader string, logger logr.Logger, httpClient *http.Client) ([]byte, error) {
	httpBasicAuthPrefix := fmt.Sprintf("%s %s", httpBasicAuthKey, authHeader)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%s %q: %w", errCreatingRequest, http.MethodPost, err)
	}
	req.Header.Set(httpAuthorizationKey, httpBasicAuthPrefix)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.New(errDoingRequest)
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
		return nil, errors.New(errReadingRequest)
	}
	return body, nil
}
