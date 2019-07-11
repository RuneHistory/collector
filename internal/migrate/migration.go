package migrate

import (
	"database/sql"
	"fmt"
	"github.com/RuneHistory/collector/internal/migrate/migrations"
	"log"
	"time"
)

type Migration interface {
	GetName() string
	Up(db *sql.DB) error
}

type MigrationRecord struct {
	ID   int
	Name string
}

func Migrate(db *sql.DB, migrations []Migration) error {
	log.Println("starting migrations")
	err := createMigrationsTable(db)
	if err != nil {
		return err
	}

	currentStage, err := getCurrentStage(db)
	if err != nil {
		return err
	}
	if currentStage.Name == "" {
		log.Printf("will run %d migrations\n", len(migrations))
		return runMigrations(db, migrations)
	}
	unstartedMigrations := make([]Migration, 0)

	foundStart := false
	for _, m := range migrations {
		if foundStart {
			unstartedMigrations = append(unstartedMigrations, m)
		}
		if m.GetName() == currentStage.Name {
			foundStart = true
		}
	}
	if len(unstartedMigrations) == 0 {
		log.Println("no migrations to run")
		return nil
	}
	log.Printf("running %d migrations\n", len(unstartedMigrations))
	return runMigrations(db, unstartedMigrations)
}

func createMigrationsTable(db *sql.DB) error {
	log.Println("creating migrations table")
	m := &migrations.CreateMigrationsTable{}
	err := m.Up(db)
	if err != nil {
		return fmt.Errorf("could not create migrations table: %s", err)
	}
	return nil
}

func runMigrations(db *sql.DB, migrations []Migration) error {
	for _, m := range migrations {
		log.Printf("running migration: %s\n", m.GetName())
		err := m.Up(db)
		if err != nil {
			log.Printf("migration %s failed\n", m.GetName())
			updateErr := setMigrationExecuted(db, m, false)
			if updateErr != nil {
				log.Printf("failed to updated migration status: %s\n", m.GetName())
				return fmt.Errorf("failed to update migration with failure: %s", updateErr)
			}
			return fmt.Errorf("migration failed: %s", err)
		}
		err = setMigrationExecuted(db, m, true)
		if err != nil {
			log.Printf("failed to updated migration status: %s\n", m.GetName())
			return fmt.Errorf("migration failed: %s", err)
		}
	}
	return nil
}

func getCurrentStage(db *sql.DB) (*MigrationRecord, error) {
	r := &MigrationRecord{}
	err := db.QueryRow("SELECT id, name FROM migrations WHERE success = 1 ORDER BY id DESC limit 1").Scan(&r.ID, &r.Name)
	if err == sql.ErrNoRows {
		return r, nil
	}
	if err != nil {
		return r, fmt.Errorf("could not find current migration stage: %v", err)
	}
	return r, nil
}

func setMigrationExecuted(db *sql.DB, m Migration, success bool) error {
	_, err := db.Exec("REPLACE INTO migrations (name, dt_executed, success) VALUES (?, ?, ?)", m.GetName(), time.Now(), success)
	return err
}
