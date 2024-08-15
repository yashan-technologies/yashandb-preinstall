package osinfoutil

import (
	"fmt"

	"preinstall/utils/executil"
	"preinstall/utils/sysutil/userutil"

	"git.yasdb.com/go/yaslog"
)

var (
	_su     = "su"
	_ulimit = "ulimit"
)

func SetSuCommand(su string) {
	_su = su
}

func SetUlimitCommand(ulimit string) {
	_ulimit = ulimit
}

func CheckUlimit(log yaslog.YasLog, username string) (string, error) {
	execer := executil.NewExecer(log)

	current, err := userutil.GetCurrentUser()
	if err != nil {
		return "", err
	}

	if current == username {
		ret, stdout, stderr := execer.Exec(_ulimit, "-a")
		if ret != 0 {
			return "", executil.GenerateError(stdout, stderr)
		}
		return stdout, nil
	}

	if userutil.IsCurrentUserRoot() {
		ret, stdout, stderr := execer.Exec(_su, username, "-c", fmt.Sprintf("%s -a", _ulimit))
		if ret != 0 {
			return "", executil.GenerateError(stdout, stderr)
		}
		return stdout, nil
	}
	return "", nil
}
