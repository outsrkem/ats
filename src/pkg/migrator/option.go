package migrator

import (
	"embed"
)

// Option configures the migrator
type Option func(*gormMigrator)

// WithTableName customizes the migration history table name
func WithTableName(name string) Option {
	return func(m *gormMigrator) {
		m.tableName = name
	}
}

// WithMigrationsDir customizes the local migration script directory
func WithMigrationsDir(dir string) Option {
	return func(m *gormMigrator) {
		if localProvider, ok := m.scriptProvider.(*LocalFileProvider); ok {
			localProvider.baseDir = dir
		}
	}
}

// WithEmbedMigrationsDir customizes the embedded migration script directory
func WithEmbedMigrationsDir(dir string, fs embed.FS) Option {
	return func(m *gormMigrator) {
		m.SetScriptProvider(NewEmbeddedScriptProvider(dir, fs))
	}
}
