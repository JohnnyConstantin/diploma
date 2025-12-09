package db

import (
	"errors"
	"fmt"

	"diploma/pkg/logger"
	"diploma/pkg/logger/message"

	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file" // file://
)

func RunMigrations(conn *SqlConnection, migrationsDir, dbName string) error {
	driver, err := migratepgx.WithInstance(conn.SqlDB, &migratepgx.Config{})
	if err != nil {
		return fmt.Errorf("error creating pgx v5 driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsDir,
		dbName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error in apply migrations: %w", err)
	}

	return nil
}

// PrepareDB запускает миграции для БД и логирует результат.
func PrepareDB(conn *SqlConnection) {
	if err := RunMigrations(conn, conn.MigrationsDir, conn.DBName); err != nil {
		logger.Log.Error(&message.LogMessage{
			Message: fmt.Sprintf("DB migrations ERROR: %s", err),
		})
		return
	}

	logger.Log.Info(&message.LogMessage{
		Message: "DB migrations applied successfully",
	})
}
