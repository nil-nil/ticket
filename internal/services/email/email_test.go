package email

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nil-nil/ticket/internal/domain"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func TestValidateSenderAddress(t *testing.T) {
	repo := &mockMailServerRepository{
		authoritativeDomains: []string{"example.com"},
	}

	server := NewServer(repo, mockCache, &mockEventBusDriver{}, func(username, password string) (domain.User, error) { return domain.User{}, nil })

	t.Run("test valid sender", func(t *testing.T) {
		err := server.ValidateSenderAddress("alan@example.com")
		assert.NoError(t, err)
	})
}

func TestValidateRecipientAddress(t *testing.T) {
	repo := &mockMailServerRepository{
		authoritativeDomains: []string{"test.com"},
		aliases: []domain.Alias{
			{User: "test", Domain: "test.com", ID: uuid.New()},
			{User: "bob", Domain: "test.com", ID: uuid.New(), DeletedAt: ptr.To(time.Now())},
		},
	}

	server := NewServer(repo, mockCache, &mockEventBusDriver{}, func(username, password string) (domain.User, error) { return domain.User{}, nil })

	table := []struct {
		description string
		email       string
		expectErr   error
	}{
		{description: "valid non-authoritative recipient", email: "alan@example.com", expectErr: nil},
		{description: "valid authoritative recipient", email: "test@test.com", expectErr: nil},
		{description: "invalid authoritative recipient", email: "fail@test.com", expectErr: ErrAliasNotFound},
		{description: "valid but deleted authoritative recipient", email: "bob@test.com", expectErr: ErrAliasNotFound},
	}

	for _, tc := range table {
		t.Run(tc.description, func(t *testing.T) {
			err := server.ValidateRecipientAddress(tc.email)
			if tc.expectErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectErr.Error())
			}
		})
	}
}

func TestReceiveData(t *testing.T) {
	t.Run("ValidMessage", func(t *testing.T) {
		message :=
			`MIME-Version: 1.0
Date: Tue, 12 Sep 2023 16:15:01 +0100
Message-ID: <1@example.com>
Subject: An example Subject
From: Bob A <bob@example.com>
To: test@test.com
Content-Type: multipart/alternative; boundary="0000000000001"

--0000000000001
Content-Type: text/plain; charset="UTF-8"

Body test data

--0000000000001
Content-Type: text/html; charset="UTF-8"
Content-Transfer-Encoding: quoted-printable

<div dir=3D"ltr">Body test data</div>

--0000000000001--
	`

		repo := &mockMailServerRepository{
			emails: make([]domain.Email, 0, 1),
			aliases: []domain.Alias{
				{ID: uuid.New(), Tenant: uuid.New(), User: "test", Domain: "test.com"},
			},
		}

		server := NewServer(repo, mockCache, &mockEventBusDriver{}, func(username, password string) (domain.User, error) { return domain.User{}, nil })

		err := server.ReceiveData(strings.NewReader(message))
		assert.NoError(t, err, "Valid Email shouldn't error")

		email := repo.emails[0]
		assert.Equal(t, "An example Subject", email.Subject, "subject should match")
		assert.True(t, email.Date.Equal(time.Date(2023, 9, 12, 15, 15, 01, 0, &time.Location{})), "Date should match header date")
	})

	t.Run("InvalidMessage", func(t *testing.T) {
		message := "lksfnlksgnlkesfn;kaef"

		repo := &mockMailServerRepository{
			emails: []domain.Email{},
		}

		server := NewServer(repo, mockCache, &mockEventBusDriver{}, func(username, password string) (domain.User, error) { return domain.User{}, nil })

		err := server.ReceiveData(strings.NewReader(message))
		assert.Error(t, err, "Invalid Email should error")
		assert.Equal(t, 0, len(repo.emails), "No email should be created on error")
	})
}
