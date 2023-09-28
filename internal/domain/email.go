package domain

import (
	"context"
	"net/mail"
	"strings"
	"time"
)

type Email struct {
	ID         uint64
	Subject    string
	Sender     string
	Recipients []string
	Date       time.Time
	Message    mail.Message
}

type EmailCreator interface {
	CreateEmail(ctx context.Context, email Email) (Email, error)
}

func CreateEmail(ctx context.Context, repo EmailCreator, msg mail.Message) (Email, error) {
	date, err := msg.Header.Date()
	if err != nil {
		date = time.Now()
	}

	subject := msg.Header.Get("Subject")
	sender := removeNames(msg.Header.Get("From"))
	recipients := strings.Split(msg.Header.Get("To"), ",")
	recipientEmails := make([]string, 0, len(recipients))
	for _, recipient := range recipients {
		recipientEmails = append(recipientEmails, removeNames(recipient))
	}
	return repo.CreateEmail(ctx, Email{Message: msg, Date: date, Subject: subject, Sender: sender, Recipients: recipientEmails})
}

func removeNames(address string) string {
	parsedAddress, err := mail.ParseAddress(address)
	if err != nil || parsedAddress == nil {
		return ""
	}
	return parsedAddress.Address
}
