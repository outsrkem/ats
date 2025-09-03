package migrator

import (
	"embed"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// Script information structure
type script struct {
	version string
	name    string
	path    string
}

// ------------------------------
// Local File Provider (for development environment)
// ------------------------------

type LocalFileProvider struct {
	baseDir string // Local script directory
}

func NewLocalFileProvider(baseDir string) *LocalFileProvider {
	return &LocalFileProvider{baseDir: baseDir}
}

func (p *LocalFileProvider) ListScripts(direction string) ([]script, error) {
	pattern := regexp.MustCompile(`^(v[0-9]+\.[0-9]+\.[0-9]+)_.+_` + direction + `\.sql$`)
	var scripts []script

	entries, err := os.ReadDir(p.baseDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if pattern.MatchString(name) {
			version := pattern.FindStringSubmatch(name)[1]
			scripts = append(scripts, script{
				version: version,
				name:    name,
				path:    filepath.Join(p.baseDir, name),
			})
		}
	}
	return scripts, nil
}

func (p *LocalFileProvider) ReadScript(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

// ------------------------------
// Embedded Resource Provider (for production environment)
// ------------------------------

type EmbeddedScriptProvider struct {
	baseDir string   // Embedded directory (e.g., "test_scripts")
	fs      embed.FS // Embedded filesystem
}

func NewEmbeddedScriptProvider(baseDir string, fs embed.FS) *EmbeddedScriptProvider {
	return &EmbeddedScriptProvider{baseDir: baseDir, fs: fs}
}

func (p *EmbeddedScriptProvider) ListScripts(direction string) ([]script, error) {
	pattern := regexp.MustCompile(`^(v[0-9]+\.[0-9]+\.[0-9]+)_.+_` + direction + `\.sql$`)
	var scripts []script

	entries, err := fs.ReadDir(p.fs, p.baseDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if pattern.MatchString(name) {
			version := pattern.FindStringSubmatch(name)[1]
			scripts = append(scripts, script{
				version: version,
				name:    name,
				path:    path.Join(p.baseDir, name),
			})
		}
	}
	return scripts, nil
}

func (p *EmbeddedScriptProvider) ReadScript(path string) (string, error) {
	content, err := fs.ReadFile(p.fs, path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}
