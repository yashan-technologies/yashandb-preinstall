package setuser

import (
	"preinstall/defines/confdef"
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
		console.OK("主用户组已存在：" + primaryGroup)
	} else {
		console.OK("创建主用户组：" + primaryGroup)
	}

	for _, group := range supplementaryGroups {
		exists, err := userutil.AddGroupIfNotExists(log.Sugar, group)
		if err != nil {
			return err
		}
		if exists {
			console.OK("附加用户组已存在：" + group)
		} else {
			console.OK("创建附加用户组：" + group)
		}
	}
	return nil
}
