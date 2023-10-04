package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AliasRepository interface {
	Find(ctx context.Context, Tenant uuid.UUID, params FindAliasParameters) (Alias, error)
	Create(ctx context.Context, Tenant uuid.UUID, user string, domain string) (Alias, error)
	Delete(ctx context.Context, Tenant uuid.UUID, ID uint64) (Alias, error)
}

type FindAliasParameters struct {
	ID     *uint64
	Domain *string
	User   *string
}

type Alias struct {
	ID        uint64
	Tenant    uuid.UUID
	User      string
	Domain    string
	DeletedAt *time.Time
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

func (s *AliasService) Find(ctx context.Context, Tenant uuid.UUID, params FindAliasParameters) (Alias, error) {
	return s.repo.Find(ctx, Tenant, params)
}

func (s *AliasService) Create(ctx context.Context, Tenant uuid.UUID, user string, domain string) (Alias, error) {
	alias, err := s.repo.Create(ctx, Tenant, user, domain)
	if err != nil {
		return Alias{}, err
	}

	return alias, nil
}

func (s *AliasService) Delete(ctx context.Context, Tenant uuid.UUID, ID uint64) (Alias, error) {
	alias, err := s.repo.Delete(ctx, Tenant, ID)
	if err != nil {
		return Alias{}, err
	}

	return alias, nil
}
