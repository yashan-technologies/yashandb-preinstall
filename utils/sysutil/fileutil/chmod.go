package fileutil

import (
	"os"
	"path/filepath"
)

// RecursiveChmod 更改指定路径下所有文件和目录的权限
func RecursiveChmod(path string, mode os.FileMode) error {
	// 使用filepath.Walk递归遍历目录
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 更改文件或目录的权限
		if err := os.Chmod(path, mode); err != nil {
			return err
		}
		return nil
	})
}

// 从底层向上递归遍历目录
func ReverseChmod(path string, mode os.FileMode) error {
	for {
		if err := os.Chmod(path, mode); err != nil {
			return err
		}
		if path == "/" {
			break
		}
		path = filepath.Dir(path)
	}
	return nil
}
