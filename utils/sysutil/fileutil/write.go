package fileutil

import (
	"io/fs"
	"os"

	"preinstall/commons/consts"

	gfs "git.yasdb.com/go/yasutil/fs"
)

const (
	DEFAULT_FILE_MODE fs.FileMode = 0644
)

func WriteFile(fname string, data []byte) error {
	return os.WriteFile(fname, data, DEFAULT_FILE_MODE)
}

func BackupFile(fname string) error {
	return gfs.CopyFile(fname, fname+consts.BakupExt())
}
