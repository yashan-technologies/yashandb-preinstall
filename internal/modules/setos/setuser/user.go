package setuser

import (
	"fmt"

	"preinstall/defines/confdef"
	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sysutil/userutil"
)

func AddUser() error {
	username := confdef.Conf().YashanDBUser.User
	homePath := confdef.Conf().YashanDBUser.Home
	primaryGroup := confdef.Conf().YashanDBUser.Group
	supplementaryGroups := confdef.Conf().YashanDBUser.AdditionalGroups
	exists, err := userutil.AddUserIfNotExists(log.Sugar, username, primaryGroup, homePath, supplementaryGroups...)
	if err != nil {
		return err
	}
	console.OK(fmt.Sprintf(i18n.T("set.user.add_to_groups"), username, supplementaryGroups))

	if !exists {
		console.OK(fmt.Sprintf(i18n.T("set.user.created"), username))
		console.Warn(i18n.T("set.user.passwd_tip"))
	} else {
		console.OK(fmt.Sprintf(i18n.T("set.user.exist"), username))
	}

	return nil
}
