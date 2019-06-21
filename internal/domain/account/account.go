package account

import "time"

type Account struct {
	ID          string
	BucketID    string
	CollectedAt time.Time
}

func NewAccount(id string, bucketId string, collectedAt time.Time) *Account {
	return &Account{
		ID:          id,
		BucketID:    bucketId,
		CollectedAt: collectedAt,
	}
}
