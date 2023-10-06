package domain

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Email struct {
	ID        uuid.UUID
	Tenant    uuid.UUID
	Subject   string
	Sender    string
	Recipient uuid.UUID
	Date      time.Time
	Message   mail.Message
}

type EmailCreator interface {
	CreateEmails(ctx context.Context, emails []Email) error
	FindAlias(ctx context.Context, user, domain string) (Alias, error)
}

func CreateEmail(ctx context.Context, repo EmailCreator, msg mail.Message) ([]Email, error) {
	date, err := msg.Header.Date()
	if err != nil {
		date = time.Now()
	}

	subject := msg.Header.Get("Subject")
	sender := removeNames(msg.Header.Get("From"))
	recipients := strings.Split(msg.Header.Get("To"), ",")
	emails := make([]Email, 0, len(recipients))
	for _, recipient := range recipients {
		address := removeNames(recipient)
		parts := strings.Split(address, "@")
		alias, err := repo.FindAlias(ctx, parts[0], parts[1])
		if errors.Is(err, ErrNotFound) {
			continue
		} else if err != nil {
			// TODO: log the error
			return nil, fmt.Errorf("Alias not available: %w", err)
		}
		emails = append(emails, Email{ID: uuid.New(), Tenant: alias.Tenant, Message: msg, Date: date, Subject: subject, Sender: sender, Recipient: alias.ID})
	}
	err = repo.CreateEmails(ctx, emails)
	if err != nil {
		return nil, fmt.Errorf("unable to create emails: %w", err)
	}
	return emails, nil
}

func removeNames(address string) string {
	parsedAddress, err := mail.ParseAddress(address)
	if err != nil || parsedAddress == nil {
		return ""
	}
	return parsedAddress.Address
}
