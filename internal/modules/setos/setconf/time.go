package setconf

import (
	"preinstall/defines/confdef"
	"preinstall/internal/modules/modulecommons/check"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sysutil/osinfoutil"
)

func SetTime() {
	if osinfoutil.IsTimeUTC8() {
		return
	}
	console.Set("设置时区")
	if err := check.CheckRootPrivilege(); err != nil {
		console.Fail("权限不足，无法设置时区，请使用root用户或者sudo执行")
		return
	}

	if err := osinfoutil.SetTimezone(log.Sugar, confdef.Conf().HostSetting.Timezone); err != nil {
		console.Fail("设置时区失败：" + err.Error())
		return
	}
	console.OK("设置时区为：" + confdef.Conf().HostSetting.Timezone)
	if err := osinfoutil.SetNtp(log.Sugar, true); err != nil {
		console.Warn("设置NTP失败[timedatectl set-ntp true]：" + err.Error())
		return
	}
}
