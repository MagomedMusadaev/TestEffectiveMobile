package migrations

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log/slog"
)

func RunMigration(db *sql.DB) error {
	const op = "internal.migrations.runMigration"

	// Инициализируем миграции
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error(op, "Ошибка при создании драйвера DB", slog.String("error", err.Error()))
		return err
	}

	// Указываем папку с миграциями
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		slog.Error(op, "Ошибка при создании миграций", slog.String("error", err.Error()))
		return err
	}

	// Применяем миграции
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error(op, "Ошибка применения миграций", slog.String("error", err.Error()))
		return err
	}

	slog.Info("Миграция успешно выполнена")
	return nil
}
