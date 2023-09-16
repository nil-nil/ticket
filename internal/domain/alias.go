package domain

import (
	"context"
	"fmt"
)

type AliasRepository interface {
	Find(context.Context, FindAliasParameters) (Alias, error)
}

type FindAliasParameters struct {
	ID     *uint64
	Domain *string
	User   *string
}

type Alias struct {
	ID     uint64
	User   string
	Domain string
}

func (a *Alias) GetEmail() string {
	return fmt.Sprintf("%s@%s", a.User, a.Domain)
}

func NewAliasService(repo AliasRepository) *AliasService {
	return &AliasService{
		repo: repo,
	}
}

type AliasService struct {
	repo AliasRepository
}

func (s *AliasService) Find(ctx context.Context, params FindAliasParameters) (Alias, error) {
	return s.repo.Find(ctx, params)
}
