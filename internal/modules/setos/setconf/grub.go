package setconf

import (
	"fmt"

	"preinstall/i18n"
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
	console.Set(i18n.T("set.grub.title"))

	if err := check.CheckRootPrivilege(); err != nil {
		console.Fail(i18n.T("set.grub.no_privilege"))
		return
	}

	kd, err := iniutil.LoadKeyData(grub_config_path, key_grub_cmdline_linux)
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("set.grub.read_fail"), err.Error()))
		return
	}

	if err := kd.Backup(); err != nil {
		console.Fail(fmt.Sprintf(i18n.T("set.grub.backup_fail"), err))
		return
	}
	console.OK(i18n.T("set.grub.backup_success"))

	if err := kd.Append(key_grub_cmdline_linux, transparent_hugepage_value, numa_value); err != nil {
		console.Fail(fmt.Sprintf(i18n.T("set.grub.update_fail"), err))
		return
	}
	console.OK(fmt.Sprintf(i18n.T("set.grub.update_success"), key_grub_cmdline_linux, []string{transparent_hugepage_value, numa_value}))

	files, err := setutil.UpdateGrubConfig(log.Sugar)
	if err != nil {
		if err == setutil.ErrMultipleGrubCfg {
			console.Fail(fmt.Sprintf(i18n.T("set.grub.multiple_cfg"), files))
			return
		}
		console.Fail(fmt.Sprintf(i18n.T("set.grub.regen_fail"), err))
		return
	}
	if len(files) > 0 {
		console.OK(fmt.Sprintf(i18n.T("set.grub.regen_success"), files[0]))
	} else {
		console.OK(i18n.T("set.grub.regen_success_default"))
	}

	if err := setutil.SetTransparentHugePageEnabledToNever(log.Sugar); err != nil {
		console.Fail(fmt.Sprintf(i18n.T("set.grub.thp_fail"), err))
		return
	}
	console.OK(i18n.T("set.grub.thp_success"))
	console.Warn(i18n.T("set.grub.reboot_tip"))
	console.Done()
}
