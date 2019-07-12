package bucket

type Repository interface {
	Get() ([]*Bucket, error)
	GetById(id string) (*Bucket, error)
	CountId(id string) (int, error)
	Create(b *Bucket) (*Bucket, error)
	Update(b *Bucket) (*Bucket, error)
	IncrementAmount(b *Bucket, amount int) error
}
