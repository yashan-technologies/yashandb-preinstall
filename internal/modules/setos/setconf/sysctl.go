package setconf

import (
	"os"
	"path"
	"strings"

	"preinstall/defines/confdef"
	"preinstall/internal/modules/modulecommons/check"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sysutil/setutil"

	"git.yasdb.com/go/yasutil/fs"
)

func UpdateSysctl() {
	console.Set("配置内核参数")
	if err := check.CheckRootPrivilege(); err != nil {
		console.Fail("权限不足，无法修改内核参数，请使用root用户或者sudo执行")
		return
	}

	if err := setutil.UpdateSysctlConf(log.Sugar, confdef.Conf().HostSetting.Sysctl...); err != nil {
		console.Fail("配置内核参数失败：" + err.Error())
		return
	}
	console.OK("备份内核参数[/etc/sysctl.conf.preinstall-bak.*]")
	console.OK("配置内核参数[/etc/sysctl.conf]")

	if err := setutil.Swapoff(log.Sugar); err != nil {
		console.Fail("关闭交换分区失败：" + err.Error())
		return
	}
	console.OK("关闭交换分区[swapoff -a]")

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
			console.Fail("创建CoreDump目录失败：" + err.Error())
			return
		}
		if err := os.Chmod(dir, 0777); err != nil {
			console.Fail("设置CoreDump目录权限失败：" + err.Error())
			return
		}
		console.OK("创建CoreDump目录[" + dir + "]")
	}

	console.Done()
}
