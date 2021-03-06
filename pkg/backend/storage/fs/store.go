package fs

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/justwatchcom/gopass/pkg/fsutil"
	"github.com/justwatchcom/gopass/pkg/out"

	"github.com/blang/semver"
)

// Store is a fs based store
type Store struct {
	path string
}

// New creates a new store
func New(dir string) *Store {
	return &Store{
		path: dir,
	}
}

// Get retrieves the named content
func (s *Store) Get(ctx context.Context, name string) ([]byte, error) {
	path := filepath.Join(s.path, filepath.Clean(name))
	out.Debug(ctx, "fs.Get(%s) - %s", name, path)
	return ioutil.ReadFile(path)
}

// Set writes the given content
func (s *Store) Set(ctx context.Context, name string, value []byte) error {
	filename := filepath.Join(s.path, filepath.Clean(name))
	filedir := filepath.Dir(filename)
	if !fsutil.IsDir(filedir) {
		if err := os.MkdirAll(filedir, 0700); err != nil {
			return err
		}
	}
	out.Debug(ctx, "fs.Set(%s) - %s", name, filepath.Join(s.path, name))
	return ioutil.WriteFile(filepath.Join(s.path, name), value, 0644)
}

// Delete removes the named entity
func (s *Store) Delete(ctx context.Context, name string) error {
	path := filepath.Join(s.path, filepath.Clean(name))
	out.Debug(ctx, "fs.Delete(%s) - %s", name, path)
	return os.Remove(path)
}

// Exists checks if the named entity exists
func (s *Store) Exists(ctx context.Context, name string) bool {
	path := filepath.Join(s.path, filepath.Clean(name))
	out.Debug(ctx, "fs.Exists(%s) - %s", name, path)
	return fsutil.IsFile(path)
}

// List returns a list of all entities
func (s *Store) List(ctx context.Context, prefix string) ([]string, error) {
	out.Debug(ctx, "fs.List(%s)", prefix)
	files := make([]string, 0, 100)
	if err := filepath.Walk(s.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") && path != s.path {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}
		if path == s.path {
			return nil
		}
		name := strings.TrimPrefix(path, s.path+string(filepath.Separator))
		if !strings.HasPrefix(name, prefix) {
			return nil
		}
		files = append(files, name)
		return nil
	}); err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

// IsDir returns true if the named entity is a directory
func (s *Store) IsDir(ctx context.Context, name string) bool {
	path := filepath.Join(s.path, filepath.Clean(name))
	out.Debug(ctx, "fs.Isdir(%s) - %s", name, path)
	return fsutil.IsDir(path)
}

// Prune removes a named directory
func (s *Store) Prune(ctx context.Context, prefix string) error {
	path := filepath.Join(s.path, filepath.Clean(prefix))
	out.Debug(ctx, "fs.Prune(%s) - %s", prefix, path)
	return os.RemoveAll(path)
}

// Name returns the name of this backend
func (s *Store) Name() string {
	return "fs"
}

// Version returns the version of this backend
func (s *Store) Version() semver.Version {
	return semver.Version{Minor: 1}
}

// Available will check if this backend is useable
func (s *Store) Available(ctx context.Context) error {
	return s.Set(ctx, ".test", []byte("test"))
}
