package setconf

import (
	"strings"

	"preinstall/defines/confdef"
	"preinstall/internal/modules/modulecommons/check"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/utils/sysutil/setutil"
)

func SetUserLimits() {
	console.Set("配置用户资源限制")

	if err := check.CheckRootPrivilege(); err != nil {
		console.Fail("权限不足，无法修改用户资源限制，请使用root用户或者sudo执行")
		return
	}

	username := confdef.Conf().YashanDBUser.User
	var limits []string
	for _, l := range confdef.Conf().YashanDBUser.Limits {
		limits = append(limits, strings.Join([]string{username, l}, " "))
	}

	if err := setutil.SetUserLimits(username, limits, true); err != nil {
		console.Fail("设置用户限制失败：" + err.Error())
		return
	}
	console.OK("配置用户资源限制[/etc/security/limits.d/" + username + ".conf]")
	console.Done()
}
