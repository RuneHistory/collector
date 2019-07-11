package migrations

import "database/sql"

type CreateMigrationsTable struct{}

func (m *CreateMigrationsTable) GetName() string {
	return "create_migrations_table"
}

func (m *CreateMigrationsTable) Up(db *sql.DB) error {
	stmt := "CREATE TABLE IF NOT EXISTS migrations (" +
		"id INT(255) AUTO_INCREMENT," +
		"name VARCHAR(255) NOT NULL," +
		"dt_executed DATETIME," +
		"success TINYINT NOT NULL," +
		"PRIMARY KEY (id)," +
		"UNIQUE KEY unique_name (name)" +
		");"
	_, err := db.Exec(stmt)
	return err
}

func (m *CreateMigrationsTable) Down(db *sql.DB) error {
	stmt := "DROP TABLE migrations;"
	_, err := db.Exec(stmt)
	return err
}
