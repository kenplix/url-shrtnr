package repository

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type FileDBConfig struct {
	Path string `mapstructure:"path"`
}

type fileDB struct {
	users *fileDBUsersRepository
	dir   string
}

func newFileDB(cfg FileDBConfig) (*fileDB, error) {
	path := cfg.Path
	if path == "" {
		path = filepath.Join("bin", "db")
	}

	err := os.MkdirAll(path, 0o700)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create %q director(y/ies)", path)
	}

	return &fileDB{dir: path}, nil
}

func (f *fileDB) close(_ context.Context) error {
	return nil
}
