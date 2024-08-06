package fileutil

import (
	"fmt"
	"os"
	"path/filepath"

	"git.yasdb.com/go/yasutil/fs"
)

func FindFile(basePath, filename string) ([]string, error) {
	if !fs.IsFileExist(basePath) {
		return nil, fmt.Errorf("base path does not exist: %s", basePath)
	}
	var files []string
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == filename {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}
