package setconf

import (
	"fmt"
	"os"
	"path"
	"strings"

	"preinstall/defines/confdef"
	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/check"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sysutil/setutil"

	"git.yasdb.com/go/yasutil/fs"
)

func UpdateSysctl() {
	console.Set(i18n.T("set.sysctl.title"))
	if err := check.CheckRootPrivilege(); err != nil {
		console.Fail(i18n.T("set.sysctl.no_privilege"))
		return
	}

	if err := setutil.UpdateSysctlConf(log.Sugar, confdef.Conf().HostSetting.Sysctl...); err != nil {
		console.Fail(fmt.Sprintf(i18n.T("set.sysctl.fail"), err.Error()))
		return
	}
	console.OK(i18n.T("set.sysctl.backup"))
	console.OK(i18n.T("set.sysctl.success"))

	if err := setutil.Swapoff(log.Sugar); err != nil {
		console.Fail(fmt.Sprintf(i18n.T("set.sysctl.swapoff_fail"), err.Error()))
		return
	}
	console.OK(i18n.T("set.sysctl.swapoff_success"))

	for _, value := range confdef.Conf().HostSetting.Sysctl {
		items := strings.Split(value, "=")
		key := strings.TrimSpace(items[0])
		if key != "kernel.core_pattern" {
			continue
		}
		value := strings.TrimSpace(strings.TrimPrefix(value, items[0]+"="))
		// 获取CoreDump的路径
		dir := path.Dir(value)
		if err := fs.Mkdir(dir); err != nil {
			console.Fail(fmt.Sprintf(i18n.T("set.sysctl.coredump_dir_fail"), err.Error()))
			return
		}
		if err := os.Chmod(dir, 0777); err != nil {
			console.Fail(fmt.Sprintf(i18n.T("set.sysctl.coredump_chmod_fail"), err.Error()))
			return
		}
		console.OK(fmt.Sprintf(i18n.T("set.sysctl.coredump_success"), dir))
	}

	console.Done()
}
