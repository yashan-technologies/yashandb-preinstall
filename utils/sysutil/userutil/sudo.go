package userutil

import (
	"strings"

	"preinstall/defines/bashdef"
	"preinstall/utils/executil"

	"git.yasdb.com/go/yaslog"
)

var _sudo = bashdef.CMD_SUDO

func SetSudoCommand(command string) {
	_sudo = command
}

func CheckSudoForUser(logger yaslog.YasLog, username string) bool {
	exec := executil.NewExecer(logger)
	// 使用sudo -l -U <username>命令来检查指定用户的sudo权限
	_, stdout, stderr := exec.Exec(_sudo, "-l", "-U", username)
	return !strings.Contains(stdout+stderr, "not allowed")
}
