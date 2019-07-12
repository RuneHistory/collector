package validate

import (
	"errors"
	"fmt"
	"github.com/RuneHistory/collector/internal/application/domain/bucket"
)

const (
	BucketIDLength = 36
)

type BucketRules interface {
	IDIsPresent(b *bucket.Bucket) error
	IDIsCorrectLength(b *bucket.Bucket) error
	IDWillBeUnique(b *bucket.Bucket) error
	IDIsUnique(b *bucket.Bucket) error
	AmountIsPositive(b *bucket.Bucket) error
	CreatedAtIsPresent(b *bucket.Bucket) error
}

func NewBucketRules(bucketRepo bucket.Repository) BucketRules {
	return &StdBucketRules{
		BucketRepo: bucketRepo,
	}
}

type StdBucketRules struct {
	BucketRepo bucket.Repository
}

func (x *StdBucketRules) IDIsPresent(b *bucket.Bucket) error {
	if b.ID == "" {
		return errors.New("id is blank")
	}
	return nil
}

func (x *StdBucketRules) IDIsCorrectLength(b *bucket.Bucket) error {
	if len(b.ID) != BucketIDLength {
		return fmt.Errorf("id %s must be exactly %d characters", b.ID, BucketIDLength)
	}
	return nil
}

func (x *StdBucketRules) IDWillBeUnique(b *bucket.Bucket) error {
	amount, err := x.BucketRepo.CountId(b.ID)
	if err != nil {
		return err
	}
	if amount != 0 {
		return fmt.Errorf("ID %s must be unique", b.ID)
	}
	return nil
}

func (x *StdBucketRules) IDIsUnique(b *bucket.Bucket) error {
	count, err := x.BucketRepo.CountId(b.ID)
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("ID %s must be unique", b.ID)
	}
	return nil
}

func (x *StdBucketRules) AmountIsPositive(b *bucket.Bucket) error {
	if b.Amount < 0 {
		return fmt.Errorf("amount %d must be >= 0", b.Amount)
	}
	return nil
}

func (x *StdBucketRules) CreatedAtIsPresent(b *bucket.Bucket) error {
	if b.CreatedAt.IsZero() {
		return errors.New("created at must be set")
	}
	return nil
}
