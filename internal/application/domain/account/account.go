package account

import "time"

type Account struct {
	ID        string
	BucketID  string
	Nickname  string
	CreatedAt time.Time
}

func NewAccount(id string, bucketId string, nickname string, createdAt time.Time) *Account {
	return &Account{
		ID:        id,
		BucketID:  bucketId,
		Nickname:  nickname,
		CreatedAt: createdAt,
	}
}
