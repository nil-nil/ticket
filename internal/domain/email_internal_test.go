package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveNames(t *testing.T) {
	table := []struct {
		description string
		address     string
		expect      string
	}{
		{description: "email only", address: "bob@test.com", expect: "bob@test.com"},
		{description: "email with name", address: "Baz <baz@test.com>", expect: "baz@test.com"},
		{description: "invalid", address: "Baz xyz", expect: ""},
	}
	for _, tc := range table {
		t.Run(tc.description, func(t *testing.T) {
			got := removeNames(tc.address)
			assert.Equal(t, tc.expect, got)
		})
	}
}
