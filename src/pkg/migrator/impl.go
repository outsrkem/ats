package migrator

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

type gormMigrator struct {
	db             *gorm.DB
	tableName      string
	scriptProvider ScriptProvider
	logger         Logger
}

// New creates a configurable migrator instance using functional options
//
// This constructor:
//   - Initializes a migrator with safe defaults
//   - Applies user-provided configuration options
//   - Provides a clean interface for customizing migration behavior
//
// Parameters:
//
//	db:    GORM database connection (required)
//	opts:  Zero or more configuration options (variadic)
//
// Default configuration:
//   - Migration table name: "auto_migration"
//   - Logger: Default logger (none)
//   - Script provider: None (must be set before use)
//
// Usage:
//
//	migrator := New(db) // Minimal configuration
//	migrator := New(db, WithTableName("custom_migrations"), WithLogger(customLogger))
//
// Important: The script provider must be configured before running migrations
func New(db *gorm.DB, opts ...Option) Migrator {
	// Initialize with safe defaults
	m := &gormMigrator{
		db:        db,
		tableName: "auto_migration", // Default table name
		logger:    NewDefaultLogger(Error),
	}

	// Apply all user-provided configuration options
	for _, opt := range opts {
		opt(m)
	}

	return m
}

// NewDefault creates a migrator instance with default configuration
// - db: GORM database connection
// - config: Optional configuration (nil for all defaults)
func NewDefault(db *gorm.DB, config *Config) Migrator {
	// Initialize with default values
	m := &gormMigrator{
		db:        db,
		tableName: "auto_migration",        // Default table name
		logger:    NewDefaultLogger(Error), // Default logger with error level
	}

	// Apply user configuration if provided
	if config != nil {
		// Override logger if specified
		if config.Logger != nil {
			m.logger = config.Logger
		}

		// Override table name if specified
		if config.TableName != "" {
			m.tableName = config.TableName
		}

		// Set script provider if specified
		if config.ScriptProvider != nil {
			m.SetScriptProvider(config.ScriptProvider)
		}
	}

	// Initialize migration history table
	if err := InitTable(db, m.tableName); err != nil {
		// Log error but continue - table may already exist
		m.logger.Errorf("Failed to initialize migration table: %v", err)
	}

	return m
}

func (m *gormMigrator) SetScriptProvider(provider ScriptProvider) {
	m.scriptProvider = provider
}

func (m *gormMigrator) Logger() Logger {
	return m.logger
}

func (m *gormMigrator) Up() error {
	scripts, err := m.listScripts("up")
	if err != nil {
		return fmt.Errorf("failed to list upgrade scripts: %w", err)
	}

	pendingScripts := m.filterPending(scripts, "up")
	if len(pendingScripts) == 0 {
		return nil
	}

	sort.Slice(pendingScripts, func(i, j int) bool {
		return compareVersions(pendingScripts[i].version, pendingScripts[j].version) < 0
	})

	for _, script := range pendingScripts {
		if err := m.executeScript(script.path, script.version, script.name, "up"); err != nil {
			return fmt.Errorf("failed to execute script %s: %w", script.name, err)
		}
	}
	return nil
}

func (m *gormMigrator) Down() error {
	lastUp, err := m.getLatestUpRecord()
	if err != nil {
		return err
	}

	// Convert up script name to down script name
	downScriptName := strings.Replace(lastUp.ScriptName, "_up.sql", "_down.sql", 1)
	if downScriptName == lastUp.ScriptName {
		return fmt.Errorf("rollback script not found for: %s", lastUp.ScriptName)
	}

	// Verify rollback script exists
	scripts, err := m.listScripts("down")
	if err != nil {
		return err
	}
	found := false
	for _, s := range scripts {
		if s.name == downScriptName {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("rollback script does not exist: %s", downScriptName)
	}

	// Execute the rollback script
	if err := m.executeScript(
		filepath.Join(m.scriptProvider.(*EmbeddedScriptProvider).baseDir, downScriptName),
		lastUp.Version, downScriptName, "down"); err != nil {
		return fmt.Errorf("failed to execute rollback script: %w", err)
	}

	// Mark the migration as rolled back (soft delete)
	isoTime := time.Now().Format(ISO8601TZ)
	return m.db.Table(m.tableName).Where("id = ?", lastUp.ID).Update("deleted_time", isoTime).Error
}

func (m *gormMigrator) HasPending() (bool, error) {
	scripts, err := m.listScripts("up")
	if err != nil {
		return false, err
	}
	pending := m.filterPending(scripts, "up")
	return len(pending) > 0, nil
}

func (m *gormMigrator) LatestVersion() (string, error) {
	var record MigrationRecord
	result := m.db.Table(m.tableName).
		Where("direction = ? AND status = ?", "up", "success").Order("executed_time DESC").First(&record)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", nil
	}
	return record.Version, result.Error
}

// executeScript executes a migration script within a database transaction
//
// This function:
//  1. Reads the script content from the provider
//  2. Splits the content into individual SQL statements
//  3. Executes all statements within a single transaction
//  4. Records the migration result (success/failure)
//
// Parameters:
//
//	scriptPath:  Path to the migration script
//	version:     Migration version (semantic version format)
//	scriptName:  Name of the migration script file
//	direction:   Migration direction ("up" or "down")
//
// Returns:
//
//	error: Detailed error if any step fails, nil on success
//
// Error handling:
//   - Ensures transaction rollback on any failure
//   - Handles recording errors separately from migration errors
//   - Returns critical error when migration succeeds but recording fails
func (m *gormMigrator) executeScript(scriptPath, version, scriptName, direction string) error {
	content, err := m.scriptProvider.ReadScript(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to read script: %w", err)
	}
	if content == "" {
		return errors.New("script content is empty")
	}

	// Split SQL statements into individual commands
	sqlStatements, err := splitSQL(content)
	if err != nil {
		return fmt.Errorf("failed to split script: %w", err)
	}
	if len(sqlStatements) == 0 {
		return errors.New("no valid SQL statements found")
	}

	// Begin database transaction
	tx := m.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Ensure rollback on panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Execute each SQL statement sequentially
	for i, stmt := range sqlStatements {
		if err := tx.Exec(stmt).Error; err != nil {
			err = fmt.Errorf("failed to execute SQL statement %d: %w", i+1, err)
			tx.Rollback()

			// Record failure and handle potential recording error
			if recordErr := m.recordResult(version, scriptName, direction, "fail", err.Error()); recordErr != nil {
				m.logger.Errorf("Failed to record migration failure: %v (original error: %v)", recordErr, err)
			}
			return err
		}
	}

	// Commit transaction if all statements succeeded
	if err := tx.Commit().Error; err != nil {
		err = fmt.Errorf("failed to commit transaction: %w", err)
		tx.Rollback()

		// Record failure and handle potential recording error
		if recordErr := m.recordResult(version, scriptName, direction, "fail", err.Error()); recordErr != nil {
			m.logger.Errorf("Failed to record migration failure: %v (original error: %v)", recordErr, err)
		}
		return err
	}

	// Record successful migration
	if err := m.recordResult(version, scriptName, direction, "success", ""); err != nil {
		// Critical error: migration succeeded but recording failed
		return fmt.Errorf("migration succeeded but failed to record result: %w", err)
	}

	return nil
}

// recordResult saves migration execution details to the history table
func (m *gormMigrator) recordResult(version, scriptName, direction, status, errMsg string) error {
	// Format as ISO 8601 with timezone offset
	isoTime := time.Now().Format(ISO8601TZ)

	return m.db.Table(m.tableName).Create(
		&MigrationRecord{
			Version:      version,
			ScriptName:   scriptName,
			Direction:    direction,
			ExecutedTime: isoTime,
			Status:       status,
			ErrorMsg:     errMsg,
		}).Error
}

// getLatestUpRecord retrieves the most recent successful 'up' migration record
func (m *gormMigrator) getLatestUpRecord() (*MigrationRecord, error) {
	var record MigrationRecord

	// Query for the latest successful 'up' migration
	result := m.db.Table(m.tableName).
		Where("direction = ? AND status = ?", "up", "success").
		Order("executed_time DESC").
		First(&record)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("no rollbackable upgrade record found")
	}

	return &record, result.Error
}

// listScripts retrieves migration scripts for the specified direction (up/down)
func (m *gormMigrator) listScripts(direction string) ([]script, error) {
	return m.scriptProvider.ListScripts(direction)
}

// filterPending identifies scripts that haven't been successfully executed
func (m *gormMigrator) filterPending(scripts []script, direction string) []script {
	var pending []script
	for _, s := range scripts {
		if !m.isExecuted(s.version, s.name, direction) {
			pending = append(pending, s)
		}
	}
	return pending
}

// isExecuted checks if a specific migration script has been successfully applied
func (m *gormMigrator) isExecuted(version, scriptName, direction string) bool {
	var count int64
	m.db.Table(m.tableName).
		Where("version = ? AND script_name = ? AND direction = ? AND status = ?",
			version, scriptName, direction, "success").
		Count(&count)
	return count > 0
}
