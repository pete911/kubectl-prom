package k8s

import (
	"fmt"
	"io"
	"k8s.io/client-go/tools/portforward"
	"net/http"
	"net/url"
)

type Forwarder struct {
	portForwarder *portforward.PortForwarder
	stopChan      chan struct{}
	host          string
	httpClient    *http.Client
}

// Get sends http GET request with specified params to port forwarded pod's path and returns
// status code and response body
func (f Forwarder) Get(path string, params url.Values) (int, []byte, error) {
	path, err := url.JoinPath(f.host, path)
	if err != nil {
		return 0, nil, fmt.Errorf("join request path: %w", err)
	}
	if len(params) != 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}

	resp, err := f.httpClient.Get(path)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("read response body: %w", err)
	}
	return resp.StatusCode, b, nil
}

// Stop stops port forwarding
func (f Forwarder) Stop() {
	f.stopChan <- struct{}{}
}
