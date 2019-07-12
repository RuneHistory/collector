package validate

import (
	"errors"
	"fmt"
	"github.com/RuneHistory/collector/internal/application/domain/account"
	"github.com/RuneHistory/collector/internal/application/domain/bucket"
)

const (
	AccountIDLength   = 36
	MaxNicknameLength = 12
)

type AccountRules interface {
	IDIsPresent(a *account.Account) error
	BucketIDIsPresent(a *account.Account) error
	IDIsCorrectLength(a *account.Account) error
	BucketIDIsCorrectLength(a *account.Account) error
	IDWillBeUnique(a *account.Account) error
	BucketIDExists(a *account.Account) error
	IDIsUnique(a *account.Account) error
	NicknameIsPresent(a *account.Account) error
	NicknameIsNotTooLong(a *account.Account) error
	NicknameIsUniqueToID(a *account.Account) error
}

func NewAccountRules(accountRepo account.Repository, bucketRepo bucket.Repository) AccountRules {
	return &StdAccountRules{
		AccountRepo: accountRepo,
		BucketRepo:  bucketRepo,
	}
}

type StdAccountRules struct {
	AccountRepo account.Repository
	BucketRepo  bucket.Repository
}

func (x *StdAccountRules) IDIsPresent(a *account.Account) error {
	if a.ID == "" {
		return errors.New("id is blank")
	}
	return nil
}

func (x *StdAccountRules) BucketIDIsPresent(a *account.Account) error {
	if a.BucketID == "" {
		return errors.New("bucket id is blank")
	}
	return nil
}

func (x *StdAccountRules) IDIsCorrectLength(a *account.Account) error {
	if len(a.ID) != AccountIDLength {
		return fmt.Errorf("id %s must be exactly %d characters", a.ID, AccountIDLength)
	}
	return nil
}

func (x *StdAccountRules) BucketIDIsCorrectLength(a *account.Account) error {
	if len(a.BucketID) != BucketIDLength {
		return fmt.Errorf("bucket id %s must be exactly %d characters", a.BucketID, BucketIDLength)
	}
	return nil
}

func (x *StdAccountRules) IDWillBeUnique(a *account.Account) error {
	amount, err := x.AccountRepo.CountId(a.ID)
	if err != nil {
		return err
	}
	if amount != 0 {
		return fmt.Errorf("ID %s must be unique", a.ID)
	}
	return nil
}

func (x *StdAccountRules) BucketIDExists(a *account.Account) error {
	amount, err := x.BucketRepo.CountId(a.BucketID)
	if err != nil {
		return err
	}
	if amount == 0 {
		return fmt.Errorf("bucket ID %s must exist", a.BucketID)
	}
	return nil
}

func (x *StdAccountRules) IDIsUnique(a *account.Account) error {
	count, err := x.AccountRepo.CountId(a.ID)
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("ID %s must be unique", a.ID)
	}
	return nil
}

func (x *StdAccountRules) NicknameIsPresent(a *account.Account) error {
	if a.Nickname == "" {
		return errors.New("nickname is blank")
	}
	return nil
}

func (x *StdAccountRules) NicknameIsNotTooLong(a *account.Account) error {
	if len(a.Nickname) > MaxNicknameLength {
		return fmt.Errorf("nickname must be no longer than %d characters", MaxNicknameLength)
	}
	return nil
}

func (x *StdAccountRules) NicknameIsUniqueToID(a *account.Account) error {
	acc, err := x.AccountRepo.GetByNicknameWithoutId(a.Nickname, a.ID)
	if err != nil {
		return err
	}
	if acc != nil {
		return fmt.Errorf("nickname %s already exists", a.Nickname)
	}
	return nil
}
