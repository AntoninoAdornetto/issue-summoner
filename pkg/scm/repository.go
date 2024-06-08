package scm

import (
	"errors"
	"os"
	"path/filepath"
)

type Repository struct {
	WorkTree string
	Dir      string
}

func NewRepository(path string) *Repository {
	return &Repository{
		WorkTree: path,
		Dir:      filepath.Join(path, ".git"),
	}
}

func FindRepository(wd string) (*Repository, error) {
	if wd == "/" {
		return nil, errors.New("unable to find local git repository")
	}

	if _, err := os.Stat(filepath.Join(wd, ".git")); err != nil {
		if os.IsNotExist(err) {
			return FindRepository(filepath.Join(wd, "../"))
		}
		return nil, err
	}

	return NewRepository(wd), nil
}
