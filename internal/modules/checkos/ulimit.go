package checkos

import (
	"fmt"
	"strings"

	"preinstall/defines/confdef"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sysutil/osinfoutil"
	"preinstall/utils/sysutil/userutil"
)

func CheckUlimit() {
	console.Check("检查数据库用户Ulimit配置")
	username := confdef.Conf().YashanDBUser.User

	if !userutil.IsUserExists(username) {
		console.Warn(fmt.Sprintf("数据库用户%s不存在，跳过检查", username))
		return
	}

	if !userutil.IsCurrentUserRoot() {
		current, err := userutil.GetCurrentUser()
		if err != nil {
			console.Fail("获取当前用户失败" + err.Error())
			return
		}
		if current != username {
			console.Fail(fmt.Sprintf("权限不足，无法检查数据库用户Ulimit配置，请使用root用户或%s用户或者sudo执行", username))
			return
		}
	}

	stdout, err := osinfoutil.CheckUlimit(log.Sugar, username)
	if err != nil {
		console.Fail("检查数据库用户Ulimit配置失败" + err.Error())
		return
	}
	console.OK(username + "用户Ulimit配置：\n" + strings.TrimSuffix(stdout, "\n"))
	console.Warn("请根据实际情况检查数据库用户Ulimit配置是否符合要求")
}
