package migrations

import "database/sql"

type CreateAccountsTable struct{}

func (m *CreateAccountsTable) GetName() string {
	return "create_accounts_table"
}

func (m *CreateAccountsTable) Up(db *sql.DB) error {
	stmt := "CREATE TABLE accounts (" +
		"id VARCHAR(36) NOT NULL," +
		"bucket_id VARCHAR(36) NOT NULL," +
		"dt_collected DATETIME," +
		"PRIMARY KEY (id)," +
		"FOREIGN KEY (bucket_id) REFERENCES buckets(id)" +
		");"
	_, err := db.Exec(stmt)
	return err
}

func (m *CreateAccountsTable) Down(db *sql.DB) error {
	stmt := "DROP TABLE accounts;"
	_, err := db.Exec(stmt)
	return err
}
