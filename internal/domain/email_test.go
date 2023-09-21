package domain_test

import (
	"context"
	"net/mail"
	"testing"
	"time"

	"github.com/nil-nil/ticket/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateEmail(t *testing.T) {
	repo := &mockCreateEmailRepository{
		emails: map[uint64]domain.Email{},
	}

	t.Run("ValidEmailNoDate", func(t *testing.T) {
		msg := mail.Message{}
		email, err := domain.CreateEmail(context.Background(), repo, msg)
		assert.NoError(t, err, "create valid email shouldn't error")
		assert.Equal(t, msg, email.Message, "message should be the same")
		assert.NotEqual(t, 0, email.ID, "ID should not be zero valued")
		gap := time.Since(email.Date)
		assert.Less(t, gap, time.Millisecond*1, "Date should be time.Now()")
		assert.Equal(t, "", email.Subject, "Missng subject header should be empty subject")
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
		email, err := domain.CreateEmail(context.Background(), repo, msg)
		assert.NoError(t, err, "create valid email shouldn't error")
		assert.Equal(t, msg, email.Message, "message should be the same")
		assert.NotEqual(t, 0, email.ID, "ID should not be zero valued")
		assert.Equal(t, "qux@example.com", email.Sender, "sender should be parsed")
		assert.Equal(t, []string{"baz@test.com", "foo@bar.com"}, email.Recipients, "recipients should be parsed")
		assert.True(t, email.Date.Equal(time.Date(2023, 9, 18, 17, 58, 07, 0, &time.Location{})), "Date should match header date")
		assert.Equal(t, "Test Message", email.Subject, "Missng subject header should be empty subject")
	})
}

type mockCreateEmailRepository struct {
	emails map[uint64]domain.Email
}

func (m *mockCreateEmailRepository) CreateEmail(ctx context.Context, email domain.Email) (domain.Email, error) {
	email.ID = nextMapKey(m.emails)
	m.emails[email.ID] = email

	return email, nil
}
