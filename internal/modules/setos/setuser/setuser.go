package setuser

import (
	"fmt"

	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/check"
	"preinstall/internal/modules/modulecommons/console"
)

func SetUser() {
	console.Set(i18n.T("set.user.title"))
	if err := check.CheckRootPrivilege(); err != nil {
		console.Fail(i18n.T("set.user.no_privilege"))
		return
	}
	if err := AddGroups(); err != nil {
		console.Fail(fmt.Sprintf(i18n.T("set.user.create_group_fail"), err.Error()))
		return
	}
	if err := AddUser(); err != nil {
		console.Fail(fmt.Sprintf(i18n.T("set.user.create_user_fail"), err.Error()))
		return
	}
	console.Done()
}
