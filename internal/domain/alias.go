package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AliasRepository interface {
	Find(ctx context.Context, Tenant uuid.UUID, params FindAliasParameters) (Alias, error)
	Create(ctx context.Context, alias Alias) error
	Delete(ctx context.Context, Tenant uuid.UUID, ID uuid.UUID) (Alias, error)
}

type FindAliasParameters struct {
	ID     *uuid.UUID
	Domain *string
	User   *string
}

type Alias struct {
	ID        uuid.UUID
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

func (s *AliasService) Create(ctx context.Context, tenant uuid.UUID, user string, domain string) (Alias, error) {
	alias := Alias{ID: uuid.New(), Tenant: tenant, User: user, Domain: domain}
	err := s.repo.Create(ctx, alias)
	if err != nil {
		return Alias{}, err
	}

	return alias, nil
}

func (s *AliasService) Delete(ctx context.Context, Tenant uuid.UUID, ID uuid.UUID) (Alias, error) {
	alias, err := s.repo.Delete(ctx, Tenant, ID)
	if err != nil {
		return Alias{}, err
	}

	return alias, nil
}
