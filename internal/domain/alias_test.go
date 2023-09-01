package domain_test

import (
	"testing"

	"github.com/nil-nil/ticket/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestAliasEmail(t *testing.T) {
	table := []struct {
		description string
		user        string
		domain      string
		expect      string
	}{
		{description: "User email", user: "bob", domain: "test.com", expect: "bob@test.com"},
	}

	for _, tc := range table {
		t.Run(tc.description, func(t *testing.T) {
			alias := domain.Alias{User: tc.user, Domain: tc.domain}
			email := alias.GetEmail()
			assert.Equal(t, tc.expect, email)
		})
	}
}
