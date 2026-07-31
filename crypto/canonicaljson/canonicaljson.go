package canonicaljson

import (
	"encoding/json"
	"fmt"
)

// canonicalizeValue replaces negative zero and floating point numbers within JSON objects
// so they can be marshaled into canonical JSON.
func canonicalizeValue(v any) (any, error) {
	switch val := v.(type) {
	case float64:
		// reject decimal numbers
		if val != float64(int64(val)) {
			return nil, fmt.Errorf("floating point numbers are not allowed in canonical JSON: %v", val)
		}
		// replace negative zero with positive zero
		if val == -0 {
			return 0, nil
		}
		return val, nil
	case []any:
		for i, elem := range val {
			canonicalizedElem, err := canonicalizeValue(elem)
			if err != nil {
				return nil, err
			}
			val[i] = canonicalizedElem
		}
		return val, nil
	case map[string]any:
		for key, elem := range val {
			canonicalizedElem, err := canonicalizeValue(elem)
			if err != nil {
				return nil, err
			}
			val[key] = canonicalizedElem
		}
		return val, nil
	default:
		return val, nil
	}
}

// Canonicalize takes a JSON string and returns the Matrix Canonical JSON representation.
func Canonicalize(input []byte) ([]byte, error) {
	// json.Marshal does the following for us:
	// keys sorted by codepoint (json.Marshal does this by default)
	// no exponents
	// numbers must be in [-(2^53)+1, 2^53-1] range
	//
	// we ourselves check for:
	// no floating points
	// no -0

	var v any
	if err := json.Unmarshal(input, &v); err != nil {
		return nil, err
	}

	v, err := canonicalizeValue(v)

	if err != nil {
		return nil, err
	}

	output, err := json.Marshal(v)

	if err != nil {
		return nil, err
	}

	return output, nil
}

// Marshal takes a Go value and returns the Matrix Canonical JSON representation.
func Marshal(v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return Canonicalize(b)
}

// IsCanonical checks if the given JSON string is in canonical form.
func IsCanonical(input []byte) (bool, error) {
	canonical, err := Canonicalize(input)
	if err != nil {
		return false, err
	}

	return string(canonical) == string(input), nil
}
