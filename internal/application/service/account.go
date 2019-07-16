package service

import (
	"fmt"
	"github.com/RuneHistory/collector/internal/application/domain"
	"github.com/RuneHistory/collector/internal/application/validate"
	"time"
)

type Account interface {
	Get() ([]*domain.Account, error)
	GetById(id string) (*domain.Account, error)
	Create(id string, bucketID string, nickname string) (*domain.Account, error)
	Update(a *domain.Account) (*domain.Account, error)
}

func NewAccountService(repo domain.AccountRepository, validator validate.AccountValidator) Account {
	return &AccountService{
		AccountRepo: repo,
		Validator:   validator,
	}
}

type AccountService struct {
	AccountRepo domain.AccountRepository
	Validator   validate.AccountValidator
}

func (s *AccountService) Get() ([]*domain.Account, error) {
	return s.AccountRepo.Get()
}

func (s *AccountService) GetById(id string) (*domain.Account, error) {
	return s.AccountRepo.GetById(id)
}

func (s *AccountService) Create(id string, bucketID string, nickname string) (*domain.Account, error) {
	now := time.Now()
	a := &domain.Account{
		ID:        id,
		BucketID:  bucketID,
		Nickname:  nickname,
		CreatedAt: now,
	}
	if err := s.Validator.NewAccount(a); err != nil {
		return nil, fmt.Errorf("unable to validate account: %s", err)
	}
	acc, err := s.AccountRepo.Create(a)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %s", err)
	}

	return acc, nil
}

func (s *AccountService) Update(a *domain.Account) (*domain.Account, error) {
	if err := s.Validator.UpdateAccount(a); err != nil {
		return nil, fmt.Errorf("unable to validate account: %s", err)
	}
	acc, err := s.AccountRepo.Update(a)
	if err != nil {
		return nil, fmt.Errorf("failed to update account: %s", err)
	}

	return acc, nil
}
