package check

import (
	"preinstall/defines/errdef"
	"preinstall/utils/sysutil/userutil"
)

func CheckRootPrivilege() error {
	if !userutil.IsCurrentUserRoot() {
		return errdef.ErrNeedRootPrivilege
	}
	return nil
}
