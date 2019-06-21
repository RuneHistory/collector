package account

type Repository interface {
	Get() ([]*Account, error)
	GetById(id string) (*Account, error)
	CountId(id string) (int, error)
	GetBySlug(slug string) (*Account, error)
	GetByNicknameWithoutId(nickname string, id string) (*Account, error)
	GetBySlugWithoutId(slug string, id string) (*Account, error)
	Create(a *Account) (*Account, error)
	Update(a *Account) (*Account, error)
}
