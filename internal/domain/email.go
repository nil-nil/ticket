package domain

import (
	"context"
	"net/mail"
	"time"
)

type Email struct {
	ID      uint64
	Subject string
	Date    time.Time
	Message mail.Message
}

type CreateEmailRepository interface {
	CreateEmail(ctx context.Context, email Email) (Email, error)
}

func CreateEmail(ctx context.Context, repo CreateEmailRepository, msg mail.Message) (Email, error) {
	date, err := msg.Header.Date()
	if err != nil {
		date = time.Now()
	}

	subject := msg.Header.Get("Subject")

	return repo.CreateEmail(ctx, Email{Message: msg, Date: date, Subject: subject})
}
