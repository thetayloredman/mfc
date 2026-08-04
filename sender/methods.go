package sender

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/thetayloredman/mfc/resolver"
)

// Sends a GET request to the given destination and path, with the given query parameters.
func (s *Sender) Get(destination, uri string) (respCode int, response map[string]any, err error) {
	xmatrix, err := s.authorizeOutboundRequest("GET", uri, destination, nil)
	if err != nil {
		return 0, nil, err
	}

	resolved, err := resolver.Resolve(destination)

	hostHeaderWithoutPort := resolved.HostHeader
	if host, _, err := net.SplitHostPort(resolved.HostHeader); err == nil {
		hostHeaderWithoutPort = host
	}

	tlsConfig := &tls.Config{
		ServerName: hostHeaderWithoutPort,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest("GET", "https://"+resolved.Endpoint+uri, nil)
	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Authorization", xmatrix)
	req.Header.Set("Content-Type", "application/json")
	req.Host = resolved.HostHeader

	fmt.Printf("sender: Sending GET %s%s (host: %s)\n", resolved.Endpoint, uri, resolved.HostHeader)
	fmt.Printf("sender: Authorization: %s\n", xmatrix)

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	var jsonResp map[string]any
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, jsonResp, nil
}

// Sends a POST request to the given destination and path, with the given body.
func (s *Sender) Post(destination, uri string, body []byte) (respCode int, response map[string]any, err error) {
	panic("TODO")
	// TODO
}

// Sends a PUT request to the given destination and path, with the given body.
func (s *Sender) Put(destination, uri string, body []byte) (respCode int, response map[string]any, err error) {
	panic("TODO")
	// TODO
}

// Sends a DELETE request to the given destination and path, with the given body.
func (s *Sender) Delete(destination, uri string, body []byte) (respCode int, response map[string]any, err error) {
	panic("TODO")
	// TODO
}

// Sends a PATCH request to the given destination and path, with the given body.
func (s *Sender) Patch(destination, uri string, body []byte) (respCode int, response map[string]any, err error) {
	panic("TODO")
	// TODO
}
