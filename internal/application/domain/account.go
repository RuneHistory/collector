package domain

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

type AccountRepository interface {
	Get() ([]*Account, error)
	GetById(id string) (*Account, error)
	GetByBucketId(id string) ([]*Account, error)
	CountId(id string) (int, error)
	GetByNicknameWithoutId(nickname string, id string) (*Account, error)
	Create(a *Account) (*Account, error)
	Update(a *Account) (*Account, error)
}
