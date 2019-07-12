package service

import (
	"fmt"
	"github.com/RuneHistory/collector/internal/application/domain/account"
	"github.com/RuneHistory/collector/internal/application/domain/validate"
	"time"
)

type Account interface {
	Get() ([]*account.Account, error)
	GetById(id string) (*account.Account, error)
	Create(id string, bucketID string, nickname string) (*account.Account, error)
	Update(a *account.Account) (*account.Account, error)
}

func NewAccountService(repo account.Repository, validator validate.AccountValidator) Account {
	return &AccountService{
		AccountRepo: repo,
		Validator:   validator,
	}
}

type AccountService struct {
	AccountRepo account.Repository
	Validator   validate.AccountValidator
}

func (s *AccountService) Get() ([]*account.Account, error) {
	return s.AccountRepo.Get()
}

func (s *AccountService) GetById(id string) (*account.Account, error) {
	return s.AccountRepo.GetById(id)
}

func (s *AccountService) Create(id string, bucketID string, nickname string) (*account.Account, error) {
	now := time.Now()
	a := &account.Account{
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

func (s *AccountService) Update(a *account.Account) (*account.Account, error) {
	if err := s.Validator.UpdateAccount(a); err != nil {
		return nil, fmt.Errorf("unable to validate account: %s", err)
	}
	acc, err := s.AccountRepo.Update(a)
	if err != nil {
		return nil, fmt.Errorf("failed to update account: %s", err)
	}

	return acc, nil
}
