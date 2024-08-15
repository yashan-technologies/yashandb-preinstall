package osinfoutil

import (
	"strings"

	"preinstall/defines/bashdef"
	"preinstall/utils/executil"

	"git.yasdb.com/go/yaslog"
)

var _systemctl = bashdef.CMD_SYSTEMCTL

func SetSystemctlCommand(command string) {
	_systemctl = command
}

func GetServiceStatus(log yaslog.YasLog, name string) (bool, string, error) {
	execer := executil.NewExecer(log)
	ret, stdout, stderr := execer.Exec(_systemctl, "is-active", name)
	if ret != 0 {
		if ret == 3 {
			return false, "inactive", nil
		}
		return false, "", executil.GenerateError(stdout, stderr)
	}
	return true, strings.TrimSpace(string(stdout)), nil
}
