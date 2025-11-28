package checkos

import (
	"fmt"
	"strings"

	"preinstall/defines/confdef"
	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sysutil/osinfoutil"
	"preinstall/utils/sysutil/userutil"
)

func CheckUlimit() {
	console.Check(i18n.T("ulimit.check"))
	username := confdef.Conf().YashanDBUser.User

	if !userutil.IsUserExists(username) {
		console.Warn(fmt.Sprintf(i18n.T("ulimit.user_not_exist"), username))
		return
	}

	if !userutil.IsCurrentUserRoot() {
		current, err := userutil.GetCurrentUser()
		if err != nil {
			console.Fail(fmt.Sprintf(i18n.T("ulimit.get_current_user_fail"), err.Error()))
			return
		}
		if current != username {
			console.Fail(fmt.Sprintf(i18n.T("ulimit.no_privilege"), username))
			return
		}
	}

	stdout, err := osinfoutil.CheckUlimit(log.Sugar, username)
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("ulimit.check_fail"), err.Error()))
		return
	}
	console.OK(fmt.Sprintf(i18n.T("ulimit.config"), username, strings.TrimSuffix(stdout, "\n")))
	console.Warn(i18n.T("ulimit.check_warn"))
}
