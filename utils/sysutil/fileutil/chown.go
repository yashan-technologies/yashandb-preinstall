package fileutil

import (
	"os"
	"path/filepath"
)

// RecursiveChown 更改指定路径下所有文件和目录的所有者
func RecursiveChown(path string, uid, gid int) error {
	// 使用filepath.Walk递归遍历目录
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 更改文件或目录的所有者
		if err := os.Chown(path, uid, gid); err != nil {
			return err
		}
		return nil
	})
}
