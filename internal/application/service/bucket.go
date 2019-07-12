package service

import (
	"fmt"
	"github.com/RuneHistory/collector/internal/application/domain/bucket"
	"github.com/RuneHistory/collector/internal/application/domain/validate"
	"github.com/satori/go.uuid"
	"sort"
	"time"
)

const MaxBucketAmount int = 10000

type Bucket interface {
	Get() ([]*bucket.Bucket, error)
	GetById(id string) (*bucket.Bucket, error)
	Create() (*bucket.Bucket, error)
	GetPriorityBucket() (*bucket.Bucket, error)
	IncrementAmount(b *bucket.Bucket) error
}

func NewBucketService(repo bucket.Repository, validator validate.BucketValidator) Bucket {
	return &BucketService{
		BucketRepo: repo,
		Validator:  validator,
	}
}

type BucketService struct {
	BucketRepo bucket.Repository
	Validator  validate.BucketValidator
}

func (s *BucketService) Get() ([]*bucket.Bucket, error) {
	return s.BucketRepo.Get()
}

func (s *BucketService) GetById(id string) (*bucket.Bucket, error) {
	return s.BucketRepo.GetById(id)
}

func (s *BucketService) Create() (*bucket.Bucket, error) {
	id := uuid.NewV4().String()
	now := time.Now()
	emptyTime := time.Time{}
	a := &bucket.Bucket{
		ID:         id,
		Amount:     0,
		CreatedAt:  now,
		StartedAt:  emptyTime,
		FinishedAt: emptyTime,
	}
	if err := s.Validator.NewBucket(a); err != nil {
		return nil, fmt.Errorf("unable to validate bucket: %s", err)
	}
	acc, err := s.BucketRepo.Create(a)
	if err != nil {
		return nil, fmt.Errorf("failed to create bucket: %s", err)
	}

	return acc, nil
}

func (s *BucketService) GetPriorityBucket() (*bucket.Bucket, error) {
	buckets, err := s.Get()
	if err != nil {
		return nil, err
	}
	if len(buckets) == 0 {
		return s.Create()
	}
	sort.Slice(buckets, func(i, j int) bool {
		return buckets[i].Amount < buckets[j].Amount
	})
	if buckets[0].Amount >= MaxBucketAmount {
		return s.Create()
	}
	return buckets[0], nil
}

func (s *BucketService) IncrementAmount(b *bucket.Bucket) error {
	err := s.BucketRepo.IncrementAmount(b, 1)
	if err != nil {
		return fmt.Errorf("failed to create bucket: %s", err)
	}
	return nil
}
