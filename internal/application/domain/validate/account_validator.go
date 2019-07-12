package validate

import "github.com/RuneHistory/collector/internal/application/domain/account"

type AccountValidator interface {
	// A set of different states/actions that may be performed.
	// Have tests against these funcs to make sure each one tests the resource
	// as expected.
	// Each one of these funcs should just use different sets of the "rules" we create in the other files/
	NewAccount(a *account.Account) error
	UpdateAccount(a *account.Account) error
}

func NewAccountValidator(accountRules AccountRules) AccountValidator {
	return &StdAccountValidator{
		accountRules: accountRules,
	}
}

type StdAccountValidator struct {
	accountRules AccountRules
}

func (v *StdAccountValidator) NewAccount(a *account.Account) error {
	if err := v.accountRules.IDIsPresent(a); err != nil {
		return err
	}
	if err := v.accountRules.IDIsCorrectLength(a); err != nil {
		return err
	}
	if err := v.accountRules.IDWillBeUnique(a); err != nil {
		return err
	}
	if err := v.accountRules.BucketIDIsPresent(a); err != nil {
		return err
	}
	if err := v.accountRules.BucketIDIsCorrectLength(a); err != nil {
		return err
	}
	if err := v.accountRules.BucketIDExists(a); err != nil {
		return err
	}
	if err := v.accountRules.NicknameIsPresent(a); err != nil {
		return err
	}
	if err := v.accountRules.NicknameIsNotTooLong(a); err != nil {
		return err
	}
	if err := v.accountRules.NicknameIsUniqueToID(a); err != nil {
		return err
	}
	return nil
}

func (v *StdAccountValidator) UpdateAccount(a *account.Account) error {
	if err := v.accountRules.IDIsPresent(a); err != nil {
		return err
	}
	if err := v.accountRules.IDIsCorrectLength(a); err != nil {
		return err
	}
	if err := v.accountRules.IDIsUnique(a); err != nil {
		return err
	}
	if err := v.accountRules.BucketIDIsPresent(a); err != nil {
		return err
	}
	if err := v.accountRules.BucketIDIsCorrectLength(a); err != nil {
		return err
	}
	if err := v.accountRules.BucketIDExists(a); err != nil {
		return err
	}
	if err := v.accountRules.NicknameIsPresent(a); err != nil {
		return err
	}
	if err := v.accountRules.NicknameIsNotTooLong(a); err != nil {
		return err
	}
	if err := v.accountRules.NicknameIsUniqueToID(a); err != nil {
		return err
	}
	return nil
}
