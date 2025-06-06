package common

import (
	"testing"
)

func TestUuidA(t *testing.T) {
	tests := []struct {
		uuid     string
		expected bool
	}{
		{"Hello123s", false},
		{"e642aba60d-874f3d9e63f77a07b0a955", false},
		{"e642aba60d874f3d9e63f77a07b0a95l", false},
		{"e642aba60d874f3d9e63f77a07b0a95", false},
		{"$%^&", false},
		{"e642aba60d874f3d9e63.77a07b0a95", false},
		{"e642aba60d874f3d9e63f77a07b0a955", true},
	}

	for _, test := range tests {
		result := CheckUuId(test.uuid)
		if result != test.expected {
			t.Errorf("WildcardMatch(%s) = %t, expected %t", test.uuid, result, test.expected)
		}
	}
}
func TestUuidB(t *testing.T) {
	tests := []struct {
		uuid     string
		expected bool
	}{
		{"Hello123s", false},
		{"e642aba60d-874f3d9e63f77a07b0a955", false},
		{"e642aba60d874f3d9e63f77a07b0a95l", false},
		{"e642aba60d874f3d9e63f77a07b0a95", false},
		{"$%^&", false},
		{"e642aba60d874f3d9e63.77a07b0a95", false},
		{"e642aba60d874f3d9e63f77a07b0a955", false},
		{"e642aba6-0d87-4f3d-9e63-f77a07b0a955", true},
	}

	for _, test := range tests {
		result := CheckStandardUUID(test.uuid)
		if result != test.expected {
			t.Errorf("WildcardMatch(%s) = %t, expected %t", test.uuid, result, test.expected)
		}
	}
}
