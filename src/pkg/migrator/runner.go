package migrator

import "fmt"

// RunWithMigrator orchestrates the database migration process
func RunWithMigrator(m Migrator) error {
	m.Logger().Infof("Checking for pending database migrations...")
	hasPending, err := m.HasPending()
	if err != nil {
		return fmt.Errorf("failed to check for pending migrations: %w", err)
	}

	if hasPending {
		m.Logger().Infof("Pending migrations detected. Applying upgrade scripts...")
		if err := m.Up(); err != nil {
			return fmt.Errorf("migration upgrade failed: %w", err)
		}
		m.Logger().Infof("Database upgrade completed successfully")
	} else {
		m.Logger().Infof("No pending migrations found")
	}

	latest, err := m.LatestVersion()
	if err != nil {
		return fmt.Errorf("failed to determine current database version: %w", err)
	}
	m.Logger().Infof("Database current version: %s", latest)
	return nil
}
