package mysql

import (
	"database/sql"
	"github.com/RuneHistory/collector/internal/application/domain/bucket"
	"github.com/go-sql-driver/mysql"
	"time"
)

func NewBucketMySQL(db *sql.DB) *BucketMySQL {
	return &BucketMySQL{
		DB: db,
	}
}

type NullableBucket struct {
	ID         string
	Amount     int
	CreatedAt  time.Time
	StartedAt  mysql.NullTime
	FinishedAt mysql.NullTime
}

type BucketMySQL struct {
	DB *sql.DB
}

func (r *BucketMySQL) Get() ([]*bucket.Bucket, error) {
	var buckets []*bucket.Bucket
	results, err := r.DB.Query("SELECT id, amount, dt_created, dt_started, dt_finished FROM buckets")
	defer func() {
		err := results.Close()
		if err != nil {
			panic(err)
		}
	}()
	if err == sql.ErrNoRows {
		return buckets, nil
	}
	if err != nil {
		return nil, err
	}
	for results.Next() {
		var nb NullableBucket
		err = results.Scan(&nb.ID, &nb.Amount, &nb.CreatedAt, &nb.StartedAt, &nb.FinishedAt)
		if err != nil {
			return nil, err
		}

		buckets = append(buckets, r.fromNullableBucket(nb))
	}
	return buckets, nil
}

func (r *BucketMySQL) GetById(id string) (*bucket.Bucket, error) {
	var nb NullableBucket
	err := r.DB.QueryRow("SELECT id, amount, dt_created, dt_started, dt_finished FROM buckets where id = ?", id).Scan(&nb.ID, &nb.Amount, &nb.CreatedAt, &nb.StartedAt, &nb.FinishedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.fromNullableBucket(nb), nil
}

func (r *BucketMySQL) CountId(id string) (int, error) {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(id) FROM buckets where id = ?", id).Scan(&count)

	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *BucketMySQL) Create(b *bucket.Bucket) (*bucket.Bucket, error) {
	nb := r.toNullableBucket(b)
	_, err := r.DB.Exec("INSERT INTO buckets (id, amount, dt_created, dt_started, dt_finished) VALUES (?, ?, ?, ?, ?)", nb.ID, nb.Amount, nb.CreatedAt, nb.StartedAt, nb.FinishedAt)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *BucketMySQL) Update(b *bucket.Bucket) (*bucket.Bucket, error) {
	nb := r.toNullableBucket(b)
	_, err := r.DB.Exec("UPDATE buckets SET dt_started = ?, dt_finished = ? WHERE id = ?", nb.StartedAt, nb.FinishedAt, nb.ID)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *BucketMySQL) IncrementAmount(b *bucket.Bucket, amount int) error {
	_, err := r.DB.Exec("UPDATE buckets SET amount = amount + ? WHERE id = ?", amount, b.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *BucketMySQL) fromNullableBucket(nb NullableBucket) *bucket.Bucket {
	startedAt := time.Time{}
	if nb.StartedAt.Valid {
		startedAt = nb.StartedAt.Time
	}
	finishedAt := time.Time{}
	if nb.FinishedAt.Valid {
		finishedAt = nb.FinishedAt.Time
	}

	return &bucket.Bucket{
		ID:         nb.ID,
		Amount:     nb.Amount,
		CreatedAt:  nb.CreatedAt,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
	}
}

func (r *BucketMySQL) toNullableBucket(b *bucket.Bucket) NullableBucket {
	startedAt := mysql.NullTime{
		Time:  b.StartedAt,
		Valid: !b.StartedAt.IsZero(),
	}
	finishedAt := mysql.NullTime{
		Time:  b.FinishedAt,
		Valid: !b.FinishedAt.IsZero(),
	}
	return NullableBucket{
		ID:         b.ID,
		Amount:     b.Amount,
		CreatedAt:  b.CreatedAt,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
	}
}
