package setutil

import (
	"preinstall/utils/executil"

	"git.yasdb.com/go/yaslog"
)

var _swapoff = "swapoff"

func SetSwapoffCommand(swapoff string) {
	_swapoff = swapoff
}

func Swapoff(log yaslog.YasLog) error {
	execer := executil.NewExecer(log)
	ret, stdout, stderr := execer.Exec(_swapoff, "-a")
	if ret != 0 {
		return executil.GenerateError(stdout, stderr)
	}
	return nil
}
