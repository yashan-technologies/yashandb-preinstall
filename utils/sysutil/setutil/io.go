package setutil

import (
	"fmt"

	"preinstall/utils/executil"

	"git.yasdb.com/go/yaslog"
)

func SetDiskQueneScheduler(log yaslog.YasLog, schedulerPath, scheduler string) error {
	execer := executil.NewExecer(log)
	cmd := fmt.Sprintf("%s %s > %s", _echo, scheduler, schedulerPath)
	ret, stdout, stderr := execer.Exec(_bash, "-c", cmd)
	if ret != 0 {
		return executil.GenerateError(stdout, stderr)
	}
	return nil
}
