package service

import (
	"fmt"
	"github.com/RuneHistory/collector/internal/application/domain"
	"github.com/RuneHistory/collector/internal/application/validate"
	"github.com/satori/go.uuid"
	"sort"
	"time"
)

const MaxBucketAmount int = 10000

type Bucket interface {
	Get() ([]*domain.Bucket, error)
	GetById(id string) (*domain.Bucket, error)
	Create() (*domain.Bucket, error)
	GetPriorityBucket() (*domain.Bucket, error)
	IncrementAmount(b *domain.Bucket) error
}

func NewBucketService(repo domain.BucketRepository, validator validate.BucketValidator) Bucket {
	return &BucketService{
		BucketRepo: repo,
		Validator:  validator,
	}
}

type BucketService struct {
	BucketRepo domain.BucketRepository
	Validator  validate.BucketValidator
}

func (s *BucketService) Get() ([]*domain.Bucket, error) {
	return s.BucketRepo.Get()
}

func (s *BucketService) GetById(id string) (*domain.Bucket, error) {
	return s.BucketRepo.GetById(id)
}

func (s *BucketService) Create() (*domain.Bucket, error) {
	id := uuid.NewV4().String()
	now := time.Now()
	emptyTime := time.Time{}
	a := &domain.Bucket{
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

func (s *BucketService) GetPriorityBucket() (*domain.Bucket, error) {
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

func (s *BucketService) IncrementAmount(b *domain.Bucket) error {
	err := s.BucketRepo.IncrementAmount(b, 1)
	if err != nil {
		return fmt.Errorf("failed to create bucket: %s", err)
	}
	return nil
}
