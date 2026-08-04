package sender

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

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
	fmt.Printf("sender: Received response: %d %+v\n", resp.StatusCode, jsonResp)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, jsonResp, nil
}

// Sends a POST request to the given destination and path, with the given body.
func (s *Sender) Post(destination, uri string, body []byte) (respCode int, response map[string]any, err error) {
	xmatrix, err := s.authorizeOutboundRequest("POST", uri, destination, body)
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

	req, err := http.NewRequest("POST", "https://"+resolved.Endpoint+uri, strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Authorization", xmatrix)
	req.Header.Set("Content-Type", "application/json")
	req.Host = resolved.HostHeader

	fmt.Printf("sender: Sending POST %s%s (host: %s)\n", resolved.Endpoint, uri, resolved.HostHeader)
	fmt.Printf("sender: Authorization: %s\n", xmatrix)

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	var jsonResp map[string]any
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	fmt.Printf("sender: Received response: %d %+v\n", resp.StatusCode, jsonResp)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, jsonResp, nil
}

// Sends a PUT request to the given destination and path, with the given body.
func (s *Sender) Put(destination, uri string, body []byte) (respCode int, response map[string]any, err error) {
	xmatrix, err := s.authorizeOutboundRequest("PUT", uri, destination, body)
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

	req, err := http.NewRequest("PUT", "https://"+resolved.Endpoint+uri, strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Authorization", xmatrix)
	req.Header.Set("Content-Type", "application/json")
	req.Host = resolved.HostHeader

	fmt.Printf("sender: Sending PUT %s%s (host: %s)\n", resolved.Endpoint, uri, resolved.HostHeader)
	fmt.Printf("sender: Authorization: %s\n", xmatrix)

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	var jsonResp map[string]any
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	fmt.Printf("sender: Received response: %d %+v\n", resp.StatusCode, jsonResp)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, jsonResp, nil
}

// Sends a DELETE request to the given destination and path, with the given body.
func (s *Sender) Delete(destination, uri string, body []byte) (respCode int, response map[string]any, err error) {
	xmatrix, err := s.authorizeOutboundRequest("DELETE", uri, destination, body)
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

	req, err := http.NewRequest("DELETE", "https://"+resolved.Endpoint+uri, strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Authorization", xmatrix)
	req.Header.Set("Content-Type", "application/json")
	req.Host = resolved.HostHeader

	fmt.Printf("sender: Sending DELETE %s%s (host: %s)\n", resolved.Endpoint, uri, resolved.HostHeader)
	fmt.Printf("sender: Authorization: %s\n", xmatrix)

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	var jsonResp map[string]any
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	fmt.Printf("sender: Received response: %d %+v\n", resp.StatusCode, jsonResp)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, jsonResp, nil
}

// Sends a PATCH request to the given destination and path, with the given body.
func (s *Sender) Patch(destination, uri string, body []byte) (respCode int, response map[string]any, err error) {
	xmatrix, err := s.authorizeOutboundRequest("PATCH", uri, destination, body)
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

	req, err := http.NewRequest("PATCH", "https://"+resolved.Endpoint+uri, strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Authorization", xmatrix)
	req.Header.Set("Content-Type", "application/json")
	req.Host = resolved.HostHeader

	fmt.Printf("sender: Sending PATCH %s%s (host: %s)\n", resolved.Endpoint, uri, resolved.HostHeader)
	fmt.Printf("sender: Authorization: %s\n", xmatrix)

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	var jsonResp map[string]any
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	fmt.Printf("sender: Received response: %d %+v\n", resp.StatusCode, jsonResp)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, jsonResp, nil
}
