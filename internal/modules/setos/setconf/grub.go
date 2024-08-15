package setconf

import (
	"fmt"

	"preinstall/internal/modules/modulecommons/check"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/iniutil"
	"preinstall/utils/sysutil/setutil"
)

const (
	key_grub_cmdline_linux     = "GRUB_CMDLINE_LINUX"
	grub_config_path           = "/etc/default/grub"
	transparent_hugepage_value = "transparent_hugepage=never"
	numa_value                 = "numa=off"
)

func DisableTransparentHugePageAndNUMA() {
	console.Set("关闭透明大页和NUMA")

	if err := check.CheckRootPrivilege(); err != nil {
		console.Fail("权限不足，无法修改相关配置，请使用root用户或者sudo执行")
		return
	}

	kd, err := iniutil.LoadKeyData(grub_config_path, key_grub_cmdline_linux)
	if err != nil {
		console.Fail("读取grub配置文件失败：" + err.Error())
		return
	}

	if err := kd.Backup(); err != nil {
		console.Fail(fmt.Sprintf("备份grub配置文件失败：%s", err))
		return
	}
	console.OK("备份grub配置文件[/ect/default/grub.preinstall-bak.*]")

	if err := kd.Append(key_grub_cmdline_linux, transparent_hugepage_value, numa_value); err != nil {
		console.Fail(fmt.Sprintf("更新grub配置文件失败：%s", err))
		return
	}
	console.OK(fmt.Sprintf("更新grub配置文件[/ect/default/grub]\n    └─ %s 配置增加 %v", key_grub_cmdline_linux, []string{transparent_hugepage_value, numa_value}))

	files, err := setutil.UpdateGrubConfig(log.Sugar)
	if err != nil {
		if err == setutil.ErrMultipleGrubCfg {
			console.Fail(fmt.Sprintf("找到多个grub.cfg：%v，请在配置文件中指定具体的grub.cfg", files))
			return
		}
		console.Fail(fmt.Sprintf("重新生成grub.cfg失败：%s", err))
		return
	}
	if len(files) > 0 {
		console.OK(fmt.Sprintf("重新生成[%s]", files[0]))
	} else {
		console.OK("重新生成grub.cfg")
	}

	if err := setutil.SetTransparentHugePageEnabledToNever(log.Sugar); err != nil {
		console.Fail(fmt.Sprintf("关闭透明大页失败：%s，您可以使用root用户执行：echo never > /sys/kernel/mm/transparent_hugepage/enabled", err))
		return
	}
	console.OK("关闭透明大页[echo never > /sys/kernel/mm/transparent_hugepage/enabled]")
	console.Done()
}
