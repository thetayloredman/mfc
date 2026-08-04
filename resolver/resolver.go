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
	ResolutionStepIpLiteral             ResolutionStep = "1 (IpLiteral)"
	ResolutionStepHostPort              ResolutionStep = "2 (HostPort)"
	ResolutionStepWellKnownIpLiteral    ResolutionStep = "3.1 (WellKnownIpLiteral)"
	ResolutionStepWellKnownHostPort     ResolutionStep = "3.2 (WellKnownHostPort)"
	ResolutionStepWellKnownSrvMatrixFed ResolutionStep = "3.3 (WellKnownSrvMatrixFed)"
	ResolutionStepWellKnownSrvMatrix    ResolutionStep = "3.4 (WellKnownSrvMatrix)"
	ResolutionStepWellKnownDefaultPort  ResolutionStep = "3.5 (WellKnownDefaultPort)"
	ResolutionStepSrvMatrixFed          ResolutionStep = "4 (SrvMatrixFed)"
	ResolutionStepSrvMatrix             ResolutionStep = "5 (SrvMatrix)"
	ResolutionStepDefaultPort           ResolutionStep = "6 (DefaultPort)"
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
	fmt.Printf("resolver: Trying to resolve %s\n", serverName)
	// 1. If the hostname is an IP literal, then that IP address should be used, together with the
	// given port number, or 8448 if no port is given.
	if isIpLiteral(serverName) {
		fmt.Printf("resolver: 1. %s is an IP literal\n", serverName)
		if host, port, err := net.SplitHostPort(serverName); err == nil {
			// if it's IPv6, wrap it in [] for the host header
			if isIp(host) && net.ParseIP(host).To4() == nil {
				host = "[" + host + "]"
			}

			fmt.Printf("resolver:    ... using host %s and explicit port %s\n", host, port)

			return &ResolvedDestination{
				Endpoint:       serverName,
				HostHeader:     host + ":" + port,
				ResolutionStep: ResolutionStepIpLiteral,
			}, nil
		}

		fmt.Printf("resolver:    ... using host %s and default port 8448\n", serverName)

		return &ResolvedDestination{
			Endpoint:       serverName + ":8448",
			HostHeader:     serverName,
			ResolutionStep: ResolutionStepIpLiteral,
		}, nil
	}

	fmt.Printf("resolver: 1. %s is not an IP literal\n", serverName)

	// 2. If the hostname is not an IP literal, and the server name includes an explicit port,
	// resolve the hostname to an IP address using CNAME, AAAA or A records.
	if _, _, err := net.SplitHostPort(serverName); err == nil {
		fmt.Printf("resolver: 2. %s includes an explicit port\n", serverName)
		fmt.Printf("resolver:    ... using host %s and explicit port\n", serverName)
		return &ResolvedDestination{
			Endpoint:       serverName,
			HostHeader:     serverName,
			ResolutionStep: ResolutionStepHostPort,
		}, nil
	}
	fmt.Printf("resolver: 2. %s does not include an explicit port\n", serverName)

	// 3. a regular HTTPS request is made to https://<hostname>/.well-known/matrix/server
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	fmt.Printf("resolver: 3. Trying https://%s/.well-known/matrix/server\n", serverName)
	resp, err := client.Get("https://" + serverName + "/.well-known/matrix/server")
	if err == nil {
		defer resp.Body.Close()

		jsonResp := struct {
			DelegatedServerName string `json:"m.server"`
		}{}

		jsonErr := json.NewDecoder(resp.Body).Decode(&jsonResp)

		if resp.StatusCode == 200 && jsonErr == nil && jsonResp.DelegatedServerName != "" {
			fmt.Printf("resolver:    Got 200 OK well-known response! Delegated to: %s\n", jsonResp.DelegatedServerName)

			// 3.1. If it's an IP literal, use that or port 8448 if no port is given.
			if isIpLiteral(jsonResp.DelegatedServerName) {
				if host, port, err := net.SplitHostPort(jsonResp.DelegatedServerName); err == nil {
					// if it's IPv6, wrap it in [] for the host header
					if isIp(host) && net.ParseIP(host).To4() == nil {
						host = "[" + host + "]"
					}

					fmt.Printf("resolver:    3.1. Delegated server name %s is an IP literal with explicit port %s\n", host, port)

					return &ResolvedDestination{
						Endpoint:       jsonResp.DelegatedServerName,
						HostHeader:     host + ":" + port,
						ResolutionStep: ResolutionStepWellKnownIpLiteral,
					}, nil
				}

				fmt.Printf("resolver:    3.1. Delegated server name %s is an IP literal with default port 8448\n", jsonResp.DelegatedServerName)

				return &ResolvedDestination{
					Endpoint:       jsonResp.DelegatedServerName + ":8448",
					HostHeader:     jsonResp.DelegatedServerName,
					ResolutionStep: ResolutionStepWellKnownIpLiteral,
				}, nil
			}

			// 3.2. If the hostname includes an explicit port, resolve the hostname to an IP address using CNAME, AAAA or A records.
			if _, _, err := net.SplitHostPort(jsonResp.DelegatedServerName); err == nil {
				fmt.Printf("resolver:    3.2. Delegated server name %s includes an explicit port\n", jsonResp.DelegatedServerName)
				return &ResolvedDestination{
					Endpoint:       jsonResp.DelegatedServerName,
					HostHeader:     jsonResp.DelegatedServerName,
					ResolutionStep: ResolutionStepWellKnownHostPort,
				}, nil
			}

			// 3.3. Try _matrix-fed._tcp.<hostname>
			fmt.Printf("resolver:          Trying SRVs...\n")
			_, addrs, err := net.LookupSRV("matrix-fed", "tcp", jsonResp.DelegatedServerName)
			if err == nil && len(addrs) > 0 {
				// remove trailing . from the target if present
				target := addrs[0].Target
				if len(target) > 0 && target[len(target)-1] == '.' {
					target = target[:len(target)-1]
				}
				fmt.Printf("resolver:    3.3.  Found SRV record for %s: %s:%d\n", jsonResp.DelegatedServerName, target, addrs[0].Port)
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
				fmt.Printf("resolver:    3.4.  Found legacy SRV record for %s: %s:%d\n", jsonResp.DelegatedServerName, target, addrs[0].Port)
				return &ResolvedDestination{
					Endpoint:       fmt.Sprintf("%s:%d", target, addrs[0].Port),
					HostHeader:     jsonResp.DelegatedServerName,
					ResolutionStep: ResolutionStepWellKnownSrvMatrix,
				}, nil
			}

			fmt.Printf("resolver:    3.5.  No SRV records found for %s, using default port 8448\n", jsonResp.DelegatedServerName)

			// 3.5. If all else fails, use the hostname with port 8448.
			return &ResolvedDestination{
				Endpoint:       jsonResp.DelegatedServerName + ":8448",
				HostHeader:     jsonResp.DelegatedServerName,
				ResolutionStep: ResolutionStepWellKnownDefaultPort,
			}, nil
		}
	}

	fmt.Printf("resolver: 3. No well-known response for %s, trying SRVs...\n", serverName)
	// 4. Try _matrix-fed._tcp.<hostname>
	_, addrs, err := net.LookupSRV("matrix-fed", "tcp", serverName)
	if err == nil && len(addrs) > 0 {
		// remove trailing . from the target if present
		target := addrs[0].Target
		if len(target) > 0 && target[len(target)-1] == '.' {
			target = target[:len(target)-1]
		}
		fmt.Printf("resolver: 4. Found SRV record for %s: %s:%d\n", serverName, target, addrs[0].Port)
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
		fmt.Printf("resolver: 5. Found legacy SRV record for %s: %s:%d\n", serverName, target, addrs[0].Port)
		return &ResolvedDestination{
			Endpoint:       fmt.Sprintf("%s:%d", target, addrs[0].Port),
			HostHeader:     serverName,
			ResolutionStep: ResolutionStepSrvMatrix,
		}, nil
	}

	fmt.Printf("resolver: 5. No SRV records found for %s, using default port 8448\n", serverName)
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
