package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

func RunMigrations(ctx context.Context, db *sqlx.DB, migrationsPath string) error {
	resolvedPath, err := resolveMigrationsPath(migrationsPath)
	if err != nil {
		return err
	}

	// Keep this table compatible with golang-migrate's default schema_migrations
	// table because this project previously used/borrowed that tooling.
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY,
			dirty BOOLEAN NOT NULL DEFAULT FALSE
		);
	`); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	files, err := os.ReadDir(resolvedPath)
	if err != nil {
		return fmt.Errorf("read migrations directory %q: %w", resolvedPath, err)
	}

	migrationFiles := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".up.sql") {
			continue
		}
		migrationFiles = append(migrationFiles, file.Name())
	}
	sort.Strings(migrationFiles)

	for _, fileName := range migrationFiles {
		version, err := migrationVersion(fileName)
		if err != nil {
			return err
		}

		applied, err := migrationApplied(ctx, db, version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		content, err := os.ReadFile(filepath.Join(resolvedPath, fileName))
		if err != nil {
			return fmt.Errorf("read migration %q: %w", fileName, err)
		}

		tx, err := db.BeginTxx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin migration %q: %w", fileName, err)
		}

		if _, err := tx.ExecContext(ctx, string(content)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("execute migration %q: %w", fileName, err)
		}

		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (version, dirty) VALUES ($1, FALSE)`, version); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %q: %w", fileName, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %q: %w", fileName, err)
		}
	}

	return nil
}

func migrationVersion(fileName string) (int64, error) {
	prefix, _, found := strings.Cut(fileName, "_")
	if !found {
		return 0, fmt.Errorf("invalid migration filename %q", fileName)
	}

	version, err := strconv.ParseInt(prefix, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid migration version %q: %w", fileName, err)
	}

	return version, nil
}

func migrationApplied(ctx context.Context, db *sqlx.DB, version int64) (bool, error) {
	var dirty bool
	if err := db.GetContext(ctx, &dirty, `SELECT COALESCE(BOOL_OR(dirty), FALSE) FROM schema_migrations`); err != nil {
		return false, fmt.Errorf("check migration dirty state: %w", err)
	}
	if dirty {
		return false, fmt.Errorf("schema_migrations is dirty; fix database migrations before starting API")
	}

	var exists bool
	if err := db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)`, version); err != nil {
		return false, fmt.Errorf("check migration %d: %w", version, err)
	}
	if exists {
		return true, nil
	}

	// golang-migrate stores only the current/latest version. If an existing
	// database says version 3 is applied, versions 1 and 2 are applied too.
	var maxApplied int64
	if err := db.GetContext(ctx, &maxApplied, `SELECT COALESCE(MAX(version), 0) FROM schema_migrations`); err != nil {
		return false, fmt.Errorf("check latest migration version: %w", err)
	}

	return maxApplied >= version, nil
}

func resolveMigrationsPath(path string) (string, error) {
	candidates := []string{path}
	if !filepath.IsAbs(path) {
		candidates = append(candidates, filepath.Join("backend", path))
	}

	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("migrations directory not found: %s", path)
}
