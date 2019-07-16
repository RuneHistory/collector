package validate

import "github.com/RuneHistory/collector/internal/application/domain"

type BucketValidator interface {
	// A set of different states/actions that may be performed.
	// Have tests against these funcs to make sure each one tests the resource
	// as expected.
	// Each one of these funcs should just use different sets of the "rules" we create in the other files/
	NewBucket(b *domain.Bucket) error
}

func NewBucketValidator(bucketRules BucketRules) BucketValidator {
	return &StdBucketValidator{
		bucketRules: bucketRules,
	}
}

type StdBucketValidator struct {
	bucketRules BucketRules
}

func (v *StdBucketValidator) NewBucket(b *domain.Bucket) error {
	if err := v.bucketRules.IDIsPresent(b); err != nil {
		return err
	}
	if err := v.bucketRules.IDIsCorrectLength(b); err != nil {
		return err
	}
	if err := v.bucketRules.IDWillBeUnique(b); err != nil {
		return err
	}
	if err := v.bucketRules.AmountIsPositive(b); err != nil {
		return err
	}
	if err := v.bucketRules.CreatedAtIsPresent(b); err != nil {
		return err
	}
	return nil
}
