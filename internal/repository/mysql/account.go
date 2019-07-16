package mysql

import (
	"database/sql"
	"github.com/RuneHistory/collector/internal/application/domain"
)

func NewAccountMySQL(db *sql.DB) *AccountMySQL {
	return &AccountMySQL{
		DB: db,
	}
}

type AccountMySQL struct {
	DB *sql.DB
}

func (r *AccountMySQL) Get() ([]*domain.Account, error) {
	var accounts []*domain.Account
	results, err := r.DB.Query("SELECT id, bucket_id, nickname, dt_created FROM accounts")
	defer func() {
		err := results.Close()
		if err != nil {
			panic(err)
		}
	}()
	if err == sql.ErrNoRows {
		return accounts, nil
	}
	if err != nil {
		return nil, err
	}
	for results.Next() {
		var a domain.Account
		err = results.Scan(&a.ID, &a.BucketID, &a.Nickname, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &a)
	}
	return accounts, nil
}

func (r *AccountMySQL) GetById(id string) (*domain.Account, error) {
	var a domain.Account
	err := r.DB.QueryRow("SELECT id, bucket_id, nickname, dt_created FROM accounts where id = ?", id).Scan(&a.ID, &a.BucketID, &a.Nickname, &a.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AccountMySQL) GetByBucketId(id string) ([]*domain.Account, error) {
	var accounts []*domain.Account
	results, err := r.DB.Query("SELECT id, bucket_id, nickname, dt_created FROM accounts WHERE bucket_id = ?", id)
	defer func() {
		err := results.Close()
		if err != nil {
			panic(err)
		}
	}()
	if err == sql.ErrNoRows {
		return accounts, nil
	}
	if err != nil {
		return nil, err
	}
	for results.Next() {
		var a domain.Account
		err = results.Scan(&a.ID, &a.BucketID, &a.Nickname, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &a)
	}
	return accounts, nil
}

func (r *AccountMySQL) CountId(id string) (int, error) {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(id) FROM accounts where id = ?", id).Scan(&count)

	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *AccountMySQL) GetByNicknameWithoutId(nickname string, id string) (*domain.Account, error) {
	var a domain.Account
	err := r.DB.QueryRow("SELECT id, bucket_id, nickname, dt_created FROM accounts where nickname = ? and id != ?", nickname, id).Scan(&a.ID, &a.BucketID, &a.Nickname, &a.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AccountMySQL) Create(a *domain.Account) (*domain.Account, error) {
	_, err := r.DB.Exec("INSERT INTO accounts (id, bucket_id, nickname, dt_created) VALUES (?, ?, ?, ?)", a.ID, a.BucketID, a.Nickname, a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *AccountMySQL) Update(a *domain.Account) (*domain.Account, error) {
	_, err := r.DB.Exec("UPDATE accounts SET nickname = ? WHERE id = ?", a.Nickname, a.ID)
	if err != nil {
		return nil, err
	}
	return a, nil
}
