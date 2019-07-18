package domain

import "time"

type Bucket struct {
	ID         string
	Amount     int
	CreatedAt  time.Time
	StartedAt  time.Time
	FinishedAt time.Time
}

func NewBucket(id string, amount int, createdAt time.Time, startedAt time.Time, finishedAt time.Time) *Bucket {
	return &Bucket{
		ID:         id,
		Amount:     amount,
		CreatedAt:  createdAt,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
	}
}

type BucketRepository interface {
	Get() ([]*Bucket, error)
	GetById(id string) (*Bucket, error)
	CountId(id string) (int, error)
	Create(b *Bucket) (*Bucket, error)
	Update(b *Bucket) (*Bucket, error)
	IncrementAmount(b *Bucket, amount int) error
}
