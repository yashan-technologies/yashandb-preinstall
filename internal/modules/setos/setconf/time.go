package setconf

import (
	"fmt"

	"preinstall/defines/confdef"
	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/check"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sysutil/osinfoutil"
)

func SetTime() {
	if osinfoutil.IsTimeUTC8() {
		return
	}
	console.Set(i18n.T("set.time.title"))
	if err := check.CheckRootPrivilege(); err != nil {
		console.Fail(i18n.T("set.time.no_privilege"))
		return
	}

	if err := osinfoutil.SetTimezone(log.Sugar, confdef.Conf().HostSetting.Timezone); err != nil {
		console.Fail(fmt.Sprintf(i18n.T("set.time.set_fail"), err.Error()))
		return
	}
	console.OK(fmt.Sprintf(i18n.T("set.time.set_success"), confdef.Conf().HostSetting.Timezone))
	if err := osinfoutil.SetNtp(log.Sugar, true); err != nil {
		console.Warn(fmt.Sprintf(i18n.T("set.time.set_fail"), err.Error()))
		return
	}
}
