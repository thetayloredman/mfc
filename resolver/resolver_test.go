package resolver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	cases := []struct {
		serverName string
		expected   ResolvedDestination
	}{
		{
			serverName: "1.2.3.4",
			expected: ResolvedDestination{
				Endpoint:       "1.2.3.4:8448",
				HostHeader:     "1.2.3.4",
				ResolutionStep: ResolutionStepIpLiteral,
			},
		},
		{
			serverName: "[::]",
			expected: ResolvedDestination{
				Endpoint:       "[::]:8448",
				HostHeader:     "[::]",
				ResolutionStep: ResolutionStepIpLiteral,
			},
		},
		{
			serverName: "1.2.3.4:8008",
			expected: ResolvedDestination{
				Endpoint:       "1.2.3.4:8008",
				HostHeader:     "1.2.3.4:8008",
				ResolutionStep: ResolutionStepIpLiteral,
			},
		},
		{
			serverName: "[::]:8008",
			expected: ResolvedDestination{
				Endpoint:       "[::]:8008",
				HostHeader:     "[::]:8008",
				ResolutionStep: ResolutionStepIpLiteral,
			},
		},
		{
			serverName: "2.s.resolvematrix.dev:7652",
			expected: ResolvedDestination{
				Endpoint:       "2.s.resolvematrix.dev:7652",
				HostHeader:     "2.s.resolvematrix.dev:7652",
				ResolutionStep: ResolutionStepHostPort,
			},
		},
		{
			serverName: "3b.s.resolvematrix.dev",
			expected: ResolvedDestination{
				Endpoint:       "wk.3b.s.resolvematrix.dev:7753",
				HostHeader:     "wk.3b.s.resolvematrix.dev:7753",
				ResolutionStep: ResolutionStepWellKnownHostPort,
			},
		},
		{
			serverName: "3c.s.resolvematrix.dev",
			expected: ResolvedDestination{
				Endpoint:       "srv.wk.3c.s.resolvematrix.dev:7754",
				HostHeader:     "wk.3c.s.resolvematrix.dev",
				ResolutionStep: ResolutionStepWellKnownSrvMatrix,
			},
		},
		{
			serverName: "3c.msc4040.s.resolvematrix.dev",
			expected: ResolvedDestination{
				Endpoint:       "srv.wk.3c.msc4040.s.resolvematrix.dev:7053",
				HostHeader:     "wk.3c.msc4040.s.resolvematrix.dev",
				ResolutionStep: ResolutionStepWellKnownSrvMatrixFed,
			},
		},
		{
			serverName: "3d.s.resolvematrix.dev",
			expected: ResolvedDestination{
				Endpoint:       "wk.3d.s.resolvematrix.dev:8448",
				HostHeader:     "wk.3d.s.resolvematrix.dev",
				ResolutionStep: ResolutionStepWellKnownDefaultPort,
			},
		},
		{
			serverName: "4.s.resolvematrix.dev",
			expected: ResolvedDestination{
				Endpoint:       "srv.4.s.resolvematrix.dev:7855",
				HostHeader:     "4.s.resolvematrix.dev",
				ResolutionStep: ResolutionStepSrvMatrix,
			},
		},
		{
			serverName: "4.msc4040.s.resolvematrix.dev",
			expected: ResolvedDestination{
				Endpoint:       "srv.4.msc4040.s.resolvematrix.dev:7054",
				HostHeader:     "4.msc4040.s.resolvematrix.dev",
				ResolutionStep: ResolutionStepSrvMatrixFed,
			},
		},
		{
			serverName: "5.s.resolvematrix.dev",
			expected: ResolvedDestination{
				Endpoint:       "5.s.resolvematrix.dev:8448",
				HostHeader:     "5.s.resolvematrix.dev",
				ResolutionStep: ResolutionStepDefaultPort,
			},
		},
	}

	for _, c := range cases {
		fmt.Printf("Testing %s\n", c.serverName)
		output, err := Resolve(c.serverName)
		assert.NoError(t, err)
		assert.Equal(t, &c.expected, output)
	}
}
