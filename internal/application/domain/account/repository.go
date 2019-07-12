package account

type Repository interface {
	Get() ([]*Account, error)
	GetById(id string) (*Account, error)
	GetByBucketId(id string) ([]*Account, error)
	CountId(id string) (int, error)
	GetByNicknameWithoutId(nickname string, id string) (*Account, error)
	Create(a *Account) (*Account, error)
	Update(a *Account) (*Account, error)
}
