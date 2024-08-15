package setuser

import (
	"fmt"

	"preinstall/defines/confdef"
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
	console.OK(fmt.Sprintf("%s用户加入附加组%v", username, supplementaryGroups))

	if !exists {
		console.OK("创建数据库用户：" + username)
		console.Warn("创建用户完成，请使用passwd命令修改数据库用户密码")
	} else {
		console.OK("数据库用户已存在：" + username)
	}

	return nil
}
