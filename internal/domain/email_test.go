package domain_test

import (
	"context"
	"net/mail"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nil-nil/ticket/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateEmail(t *testing.T) {
	baz := domain.Alias{ID: uuid.New(), User: "baz", Domain: "test.com", Tenant: uuid.New()}
	repo := &mockCreateEmailRepository{
		emails: make([]domain.Email, 0, 512),
		aliases: []domain.Alias{
			baz,
		},
	}

	t.Run("ValidEmailNoDate", func(t *testing.T) {
		msg := mail.Message{
			Header: mail.Header{
				"To":   {"Baz <baz@test.com>, foo@bar.com"},
				"From": {"Qux <qux@example.com>"},
			},
		}
		emails, err := domain.CreateEmail(context.Background(), repo, msg)
		assert.NoError(t, err, "create valid email shouldn't error")
		assert.Len(t, emails, 1, "expected 1 email to be created")
		email := emails[0]
		assert.Equal(t, msg, email.Message, "message should be the same")
		assert.NotEqual(t, uuid.Nil, email.ID, "ID should not be zero valued")
		gap := time.Since(email.Date)
		assert.Less(t, gap, time.Millisecond*1, "Date should be time.Now()")
		assert.Equal(t, "", email.Subject, "Missing subject header should be empty subject")
		assert.Equal(t, baz.ID, email.Recipient, "Expected correct recipient ID")
		assert.Equal(t, baz.Tenant, email.Tenant, "Expected correct Tenant ID")
	})

	t.Run("ValidEmail", func(t *testing.T) {
		msg := mail.Message{
			Header: mail.Header{
				"Date":    {"Mon, 18 Sep 2023 17:58:07 +0000 (UTC)"},
				"Subject": {"Test Message"},
				"To":      {"Baz <baz@test.com>, foo@bar.com"},
				"From":    {"Qux <qux@example.com>"},
			},
		}
		emails, err := domain.CreateEmail(context.Background(), repo, msg)
		assert.NoError(t, err, "create valid email shouldn't error")
		assert.Len(t, emails, 1, "expected 1 email to be created")
		email := emails[0]
		assert.Equal(t, msg, email.Message, "message should be the same")
		assert.NotEqual(t, 0, email.ID, "ID should not be zero valued")
		assert.Equal(t, "qux@example.com", email.Sender, "sender should be parsed")
		assert.True(t, email.Date.Equal(time.Date(2023, 9, 18, 17, 58, 07, 0, &time.Location{})), "Date should match header date")
		assert.Equal(t, "Test Message", email.Subject, "Expected correct subject")
		assert.Equal(t, baz.ID, email.Recipient, "Expected correct recipient ID")
		assert.Equal(t, baz.Tenant, email.Tenant, "Expected correct Tenant ID")
	})
}

type mockCreateEmailRepository struct {
	emails  []domain.Email
	aliases []domain.Alias
}

func (m *mockCreateEmailRepository) CreateEmails(ctx context.Context, emails []domain.Email) error {
	m.emails = append(m.emails, emails...)

	return nil
}

func (m *mockCreateEmailRepository) FindAlias(ctx context.Context, user, dnsDomain string) (domain.Alias, error) {
	for _, alias := range m.aliases {
		if alias.User == user && alias.Domain == dnsDomain {
			return alias, nil
		}
	}

	return domain.Alias{}, domain.ErrNotFound
}
