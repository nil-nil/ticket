package mail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSenderAddress(t *testing.T) {
	server := server{
		blockedSenders: map[string]struct{}{"bob@example.com": {}},
	}
	t.Run("test valid sender", func(t *testing.T) {
		err := server.ValidateSenderAddress("alan@example.com")
		assert.NoError(t, err)
	})

	t.Run("test blocked sender", func(t *testing.T) {
		err := server.ValidateSenderAddress("bob@example.com")
		assert.EqualError(t, err, ErrBlockedSender.Error())
	})
}

func TestValidateRecipientAddress(t *testing.T) {
	server := server{}
	t.Run("test valid recipient", func(t *testing.T) {
		err := server.ValidateRecipientAddress("alan@example.com")
		assert.NoError(t, err)
	})
}

func TestGetUserAndDomainParts(t *testing.T) {
	cases := []struct {
		description  string
		address      string
		expectUser   string
		expectDomain string
		expectErr    error
	}{
		{description: "valid email", address: "test@test.com", expectUser: "test", expectDomain: "test.com", expectErr: nil},
		{description: "missing domain", address: "bob@", expectUser: "", expectDomain: "", expectErr: ErrInvalidEmailAddress},
		{description: "missing user", address: "@nob", expectUser: "", expectDomain: "", expectErr: ErrInvalidEmailAddress},
		{description: "not even an @ sign", address: "gary", expectUser: "", expectDomain: "", expectErr: ErrInvalidEmailAddress},
	}

	for _, testCase := range cases {
		t.Run(testCase.description, func(t *testing.T) {
			user, domain, err := getUserAndDomainParts(testCase.address)
			assert.Equal(t, testCase.expectUser, user, "user doesn't match")
			assert.Equal(t, testCase.expectDomain, domain, "domain doesn't match")
			if testCase.expectErr != nil {
				assert.EqualError(t, err, testCase.expectErr.Error(), "err doesn't match")
			}
		})
	}
}
