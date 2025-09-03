package migrator

// Migrator core migration interface
type Migrator interface {
	Up() error
	Down() error
	HasPending() (bool, error)
	LatestVersion() (string, error)
	SetScriptProvider(provider ScriptProvider)
	Logger() Logger
}

// ScriptProvider abstracts script loading from different sources
type ScriptProvider interface {
	ListScripts(direction string) ([]script, error)
	ReadScript(path string) (string, error)
}

// Config holds migrator configuration options
type Config struct {
	Logger         Logger         // Logger for migration operations
	TableName      string         // Migration history table name
	ScriptProvider ScriptProvider // Source for migration scripts
}
