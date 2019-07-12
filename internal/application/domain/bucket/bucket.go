package bucket

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
