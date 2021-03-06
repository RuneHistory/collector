package migrations

import "database/sql"

type CreateBucketsTable struct{}

func (m *CreateBucketsTable) GetName() string {
	return "create_buckets_table"
}

func (m *CreateBucketsTable) Up(db *sql.DB) error {
	stmt := "CREATE TABLE buckets (" +
		"id VARCHAR(36) NOT NULL," +
		"amount int UNSIGNED NOT NULL," +
		"dt_created DATETIME NOT NULL," +
		"dt_started DATETIME," +
		"dt_finished DATETIME," +
		"PRIMARY KEY (id)" +
		");"
	_, err := db.Exec(stmt)
	return err
}

func (m *CreateBucketsTable) Down(db *sql.DB) error {
	stmt := "DROP TABLE buckets;"
	_, err := db.Exec(stmt)
	return err
}
