package setuser

import (
	"preinstall/internal/modules/modulecommons/check"
	"preinstall/internal/modules/modulecommons/console"
)

func SetUser() {
	console.Set("配置数据库用户")
	if err := check.CheckRootPrivilege(); err != nil {
		console.Fail("权限不足，无法配置用户，请使用root用户或者sudo执行")
		return
	}
	if err := AddGroups(); err != nil {
		console.Fail("创建用户组失败：" + err.Error())
		return
	}
	if err := AddUser(); err != nil {
		console.Fail("创建数据库用户失败：" + err.Error())
		return
	}
	console.Done()
}
