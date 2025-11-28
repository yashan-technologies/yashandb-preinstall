package setconf

import (
	"fmt"
	"strings"

	"preinstall/defines/confdef"
	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/check"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/utils/sysutil/setutil"
)

func SetUserLimits() {
	console.Set(i18n.T("set.limit.title"))

	if err := check.CheckRootPrivilege(); err != nil {
		console.Fail(i18n.T("set.limit.no_privilege"))
		return
	}

	username := confdef.Conf().YashanDBUser.User
	var limits []string
	for _, l := range confdef.Conf().YashanDBUser.Limits {
		limits = append(limits, strings.Join([]string{username, l}, " "))
	}

	if err := setutil.SetUserLimits(username, limits, true); err != nil {
		console.Fail(fmt.Sprintf(i18n.T("set.limit.fail"), err.Error()))
		return
	}
	console.OK(fmt.Sprintf(i18n.T("set.limit.success"), username))
	console.Done()
}
