package canonicaljson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanonicalize(t *testing.T) {
	// table based test
	cases := []struct {
		input    []byte
		expected []byte
	}{
		{
			input:    []byte(`{}`),
			expected: []byte(`{}`),
		},
		{
			input:    []byte(`{ "one": 1, "two": "Two" }`),
			expected: []byte(`{"one":1,"two":"Two"}`),
		},
		{
			input:    []byte(`{ "b": "2", "a": "1" }`),
			expected: []byte(`{"a":"1","b":"2"}`),
		},
		{
			input: []byte(`
				{
					"auth": {
						"success": true,
						"mxid": "@john.doe:example.com",
						"profile": {
							"display_name": "John Doe",
							"three_pids": [
								{
									"medium": "email",
									"address": "john.doe@example.org"
								},
								{
									"medium": "msisdn",
									"address": "123456789"
								}
							]
						}
					}
				}
			`),
			expected: []byte(`{"auth":{"mxid":"@john.doe:example.com","profile":{"display_name":"John Doe","three_pids":[{"address":"john.doe@example.org","medium":"email"},{"address":"123456789","medium":"msisdn"}]},"success":true}}`),
		},
		{
			input:    []byte(`{ "a": "日本語" }`),
			expected: []byte(`{"a":"日本語"}`),
		},
		{
			input:    []byte(`{ "a": "\u65E5" }`),
			expected: []byte(`{"a":"日"}`),
		},
		{
			input:    []byte(`{ "a": null }`),
			expected: []byte(`{"a":null}`),
		},
		{
			input:    []byte(`{ "a": -0, "b": 1e10 }`),
			expected: []byte(`{"a":0,"b":10000000000}`),
		},
	}

	for _, c := range cases {
		output, err := Canonicalize(c.input)
		assert.NoError(t, err)
		assert.Equal(t, c.expected, output)

		isCanonical, err := IsCanonical(c.expected)
		assert.NoError(t, err)
		assert.True(t, isCanonical)
	}
}

func TestCanonicalizeInvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{ "a": 1, "b": }`)
	output, err := Canonicalize(invalidJSON)
	assert.Error(t, err)
	assert.Nil(t, output)
}

func TestCanonicalizeFloatingPoint(t *testing.T) {
	floatingPointJSON := []byte(`{ "a": 1.23 }`)
	output, err := Canonicalize(floatingPointJSON)
	assert.Error(t, err)
	assert.Nil(t, output)
}
