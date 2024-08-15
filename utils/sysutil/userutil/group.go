package userutil

import (
	"os/user"

	"preinstall/utils/executil"

	"git.yasdb.com/go/yaslog"
)

func IsGroupExists(groupname string) bool {
	_, err := user.LookupGroup(groupname)
	return err == nil
}

func GetGroupByID(gid string) (*user.Group, error) {
	return user.LookupGroupId(gid)
}

func AddGroupIfNotExists(log yaslog.YasLog, groupname string) (bool, error) {
	if IsGroupExists(groupname) {
		return true, nil
	}
	execer := executil.NewExecer(log)
	ret, stdout, stderr := execer.Exec(_groupaddCommand, groupname)
	if ret != 0 {
		return false, executil.GenerateError(stdout, stderr)
	}
	return false, nil
}
