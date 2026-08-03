package resolver

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

type ResolutionStep string

const (
	ResolutionStepIpLiteral             ResolutionStep = "1"
	ResolutionStepHostPort              ResolutionStep = "2"
	ResolutionStepWellKnownIpLiteral    ResolutionStep = "3.1"
	ResolutionStepWellKnownHostPort     ResolutionStep = "3.2"
	ResolutionStepWellKnownSrvMatrixFed ResolutionStep = "3.3"
	ResolutionStepWellKnownSrvMatrix    ResolutionStep = "3.4"
	ResolutionStepWellKnownDefaultPort  ResolutionStep = "3.5"
	ResolutionStepSrvMatrixFed          ResolutionStep = "4"
	ResolutionStepSrvMatrix             ResolutionStep = "5"
	ResolutionStepDefaultPort           ResolutionStep = "6"
)

type ResolvedDestination struct {
	// The IP/hostname plus port to connect to.
	Endpoint string
	// The SNI/Host header to use when connecting.
	HostHeader string
	// The step that was used to resolve this destination.
	ResolutionStep ResolutionStep
}

func Resolve(serverName string) (*ResolvedDestination, error) {
	// 1. If the hostname is an IP literal, then that IP address should be used, together with the
	// given port number, or 8448 if no port is given.
	if isIpLiteral(serverName) {
		if host, port, err := net.SplitHostPort(serverName); err == nil {
			// if it's IPv6, wrap it in [] for the host header
			if isIp(host) && net.ParseIP(host).To4() == nil {
				host = "[" + host + "]"
			}

			return &ResolvedDestination{
				Endpoint:       serverName,
				HostHeader:     host + ":" + port,
				ResolutionStep: ResolutionStepIpLiteral,
			}, nil
		}

		return &ResolvedDestination{
			Endpoint:       serverName + ":8448",
			HostHeader:     serverName,
			ResolutionStep: ResolutionStepIpLiteral,
		}, nil
	}

	// 2. If the hostname is not an IP literal, and the server name includes an explicit port,
	// resolve the hostname to an IP address using CNAME, AAAA or A records.
	if _, _, err := net.SplitHostPort(serverName); err == nil {
		return &ResolvedDestination{
			Endpoint:       serverName,
			HostHeader:     serverName,
			ResolutionStep: ResolutionStepHostPort,
		}, nil
	}

	// 3. a regular HTTPS request is made to https://<hostname>/.well-known/matrix/server
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get("https://" + serverName + "/.well-known/matrix/server")
	if err == nil {
		defer resp.Body.Close()

		jsonResp := struct {
			DelegatedServerName string `json:"m.server"`
		}{}

		jsonErr := json.NewDecoder(resp.Body).Decode(&jsonResp)

		if resp.StatusCode == 200 && jsonErr == nil && jsonResp.DelegatedServerName != "" {
			// 3.1. If it's an IP literal, use that or port 8448 if no port is given.
			if isIpLiteral(jsonResp.DelegatedServerName) {
				if host, port, err := net.SplitHostPort(jsonResp.DelegatedServerName); err == nil {
					// if it's IPv6, wrap it in [] for the host header
					if isIp(host) && net.ParseIP(host).To4() == nil {
						host = "[" + host + "]"
					}

					return &ResolvedDestination{
						Endpoint:       jsonResp.DelegatedServerName,
						HostHeader:     host + ":" + port,
						ResolutionStep: ResolutionStepWellKnownIpLiteral,
					}, nil
				}

				return &ResolvedDestination{
					Endpoint:       jsonResp.DelegatedServerName + ":8448",
					HostHeader:     jsonResp.DelegatedServerName,
					ResolutionStep: ResolutionStepWellKnownIpLiteral,
				}, nil
			}

			// 3.2. If the hostname includes an explicit port, resolve the hostname to an IP address using CNAME, AAAA or A records.
			if _, _, err := net.SplitHostPort(jsonResp.DelegatedServerName); err == nil {
				return &ResolvedDestination{
					Endpoint:       jsonResp.DelegatedServerName,
					HostHeader:     jsonResp.DelegatedServerName,
					ResolutionStep: ResolutionStepWellKnownHostPort,
				}, nil
			}

			// 3.3. Try _matrix-fed._tcp.<hostname>
			_, addrs, err := net.LookupSRV("matrix-fed", "tcp", jsonResp.DelegatedServerName)
			if err == nil && len(addrs) > 0 {
				// remove trailing . from the target if present
				target := addrs[0].Target
				if len(target) > 0 && target[len(target)-1] == '.' {
					target = target[:len(target)-1]
				}
				return &ResolvedDestination{
					Endpoint:       fmt.Sprintf("%s:%d", target, addrs[0].Port),
					HostHeader:     jsonResp.DelegatedServerName,
					ResolutionStep: ResolutionStepWellKnownSrvMatrixFed,
				}, nil
			}

			// 3.4. Legacy _matrix._tcp.<hostname>
			_, addrs, err = net.LookupSRV("matrix", "tcp", jsonResp.DelegatedServerName)
			if err == nil && len(addrs) > 0 {
				// remove trailing . from the target if present
				target := addrs[0].Target
				if len(target) > 0 && target[len(target)-1] == '.' {
					target = target[:len(target)-1]
				}
				return &ResolvedDestination{
					Endpoint:       fmt.Sprintf("%s:%d", target, addrs[0].Port),
					HostHeader:     jsonResp.DelegatedServerName,
					ResolutionStep: ResolutionStepWellKnownSrvMatrix,
				}, nil
			}

			// 3.5. If all else fails, use the hostname with port 8448.
			return &ResolvedDestination{
				Endpoint:       jsonResp.DelegatedServerName + ":8448",
				HostHeader:     jsonResp.DelegatedServerName,
				ResolutionStep: ResolutionStepWellKnownDefaultPort,
			}, nil
		}
	}

	// 4. Try _matrix-fed._tcp.<hostname>
	_, addrs, err := net.LookupSRV("matrix-fed", "tcp", serverName)
	if err == nil && len(addrs) > 0 {
		// remove trailing . from the target if present
		target := addrs[0].Target
		if len(target) > 0 && target[len(target)-1] == '.' {
			target = target[:len(target)-1]
		}
		return &ResolvedDestination{
			Endpoint:       fmt.Sprintf("%s:%d", target, addrs[0].Port),
			HostHeader:     serverName,
			ResolutionStep: ResolutionStepSrvMatrixFed,
		}, nil
	}

	// 5. Legacy _matrix._tcp.<hostname>
	_, addrs, err = net.LookupSRV("matrix", "tcp", serverName)
	if err == nil && len(addrs) > 0 {
		// remove trailing . from the target if present
		target := addrs[0].Target
		if len(target) > 0 && target[len(target)-1] == '.' {
			target = target[:len(target)-1]
		}
		return &ResolvedDestination{
			Endpoint:       fmt.Sprintf("%s:%d", target, addrs[0].Port),
			HostHeader:     serverName,
			ResolutionStep: ResolutionStepSrvMatrix,
		}, nil
	}

	// 6. If all else fails, use the hostname with port 8448.
	return &ResolvedDestination{
		Endpoint:       serverName + ":8448",
		HostHeader:     serverName,
		ResolutionStep: ResolutionStepDefaultPort,
	}, nil
}

func isIp(hostname string) bool {
	// prune [] from ipv6 literals
	if len(hostname) > 0 && hostname[0] == '[' && hostname[len(hostname)-1] == ']' {
		hostname = hostname[1 : len(hostname)-1]
	}

	return net.ParseIP(hostname) != nil
}

func isIpLiteral(hostname string) bool {
	if isIp(hostname) {
		return true
	}

	if host, _, err := net.SplitHostPort(hostname); err == nil {
		return isIp(host)
	}

	return false
}
