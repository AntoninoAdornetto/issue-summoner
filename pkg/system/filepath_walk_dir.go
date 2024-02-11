package system

import (
	"io/fs"
	"path/filepath"
)

type FilePathDirWalker struct{}

func (dw FilePathDirWalker) WalkDir(root string, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(root, fn)
}
