package setuser

import (
	"fmt"

	"preinstall/defines/confdef"
	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sysutil/userutil"
)

func AddGroups() error {
	primaryGroup := confdef.Conf().YashanDBUser.Group
	supplementaryGroups := confdef.Conf().YashanDBUser.AdditionalGroups
	exists, err := userutil.AddGroupIfNotExists(log.Sugar, primaryGroup)
	if err != nil {
		return err
	}
	if exists {
		console.OK(fmt.Sprintf(i18n.T("set.user.primary_group_exist"), primaryGroup))
	} else {
		console.OK(fmt.Sprintf(i18n.T("set.user.primary_group_created"), primaryGroup))
	}

	for _, group := range supplementaryGroups {
		exists, err := userutil.AddGroupIfNotExists(log.Sugar, group)
		if err != nil {
			return err
		}
		if exists {
			console.OK(fmt.Sprintf(i18n.T("set.user.supplementary_group_exist"), group))
		} else {
			console.OK(fmt.Sprintf(i18n.T("set.user.supplementary_group_created"), group))
		}
	}
	return nil
}
