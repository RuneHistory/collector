package validate

import (
	"errors"
	"fmt"
	"github.com/RuneHistory/collector/internal/application/domain"
)

const (
	AccountIDLength   = 36
	MaxNicknameLength = 12
)

type AccountRules interface {
	IDIsPresent(a *domain.Account) error
	BucketIDIsPresent(a *domain.Account) error
	IDIsCorrectLength(a *domain.Account) error
	BucketIDIsCorrectLength(a *domain.Account) error
	IDWillBeUnique(a *domain.Account) error
	BucketIDExists(a *domain.Account) error
	IDIsUnique(a *domain.Account) error
	NicknameIsPresent(a *domain.Account) error
	NicknameIsNotTooLong(a *domain.Account) error
	NicknameIsUniqueToID(a *domain.Account) error
}

func NewAccountRules(accountRepo domain.AccountRepository, bucketRepo domain.BucketRepository) AccountRules {
	return &StdAccountRules{
		AccountRepo: accountRepo,
		BucketRepo:  bucketRepo,
	}
}

type StdAccountRules struct {
	AccountRepo domain.AccountRepository
	BucketRepo  domain.BucketRepository
}

func (x *StdAccountRules) IDIsPresent(a *domain.Account) error {
	if a.ID == "" {
		return errors.New("id is blank")
	}
	return nil
}

func (x *StdAccountRules) BucketIDIsPresent(a *domain.Account) error {
	if a.BucketID == "" {
		return errors.New("bucket id is blank")
	}
	return nil
}

func (x *StdAccountRules) IDIsCorrectLength(a *domain.Account) error {
	if len(a.ID) != AccountIDLength {
		return fmt.Errorf("id %s must be exactly %d characters", a.ID, AccountIDLength)
	}
	return nil
}

func (x *StdAccountRules) BucketIDIsCorrectLength(a *domain.Account) error {
	if len(a.BucketID) != BucketIDLength {
		return fmt.Errorf("bucket id %s must be exactly %d characters", a.BucketID, BucketIDLength)
	}
	return nil
}

func (x *StdAccountRules) IDWillBeUnique(a *domain.Account) error {
	amount, err := x.AccountRepo.CountId(a.ID)
	if err != nil {
		return err
	}
	if amount != 0 {
		return fmt.Errorf("ID %s must be unique", a.ID)
	}
	return nil
}

func (x *StdAccountRules) BucketIDExists(a *domain.Account) error {
	amount, err := x.BucketRepo.CountId(a.BucketID)
	if err != nil {
		return err
	}
	if amount == 0 {
		return fmt.Errorf("bucket ID %s must exist", a.BucketID)
	}
	return nil
}

func (x *StdAccountRules) IDIsUnique(a *domain.Account) error {
	count, err := x.AccountRepo.CountId(a.ID)
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("ID %s must be unique", a.ID)
	}
	return nil
}

func (x *StdAccountRules) NicknameIsPresent(a *domain.Account) error {
	if a.Nickname == "" {
		return errors.New("nickname is blank")
	}
	return nil
}

func (x *StdAccountRules) NicknameIsNotTooLong(a *domain.Account) error {
	if len(a.Nickname) > MaxNicknameLength {
		return fmt.Errorf("nickname must be no longer than %d characters", MaxNicknameLength)
	}
	return nil
}

func (x *StdAccountRules) NicknameIsUniqueToID(a *domain.Account) error {
	acc, err := x.AccountRepo.GetByNicknameWithoutId(a.Nickname, a.ID)
	if err != nil {
		return err
	}
	if acc != nil {
		return fmt.Errorf("nickname %s already exists", a.Nickname)
	}
	return nil
}
