package checkos

import (
	"fmt"
	"strings"

	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sysutil/osinfoutil"
	"preinstall/utils/sysutil/setutil"
)

func CheckSysctlConf() {
	console.Check(i18n.T("sysctl.check"))

	stdout, err := setutil.ReadSysctlConf(log.Sugar)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "permission denied") {
			console.Fail(i18n.T("sysctl.no_privilege"))
			return
		}
		console.Fail(fmt.Sprintf(i18n.T("sysctl.read_fail"), err.Error()))
		return
	}

	console.OK(fmt.Sprintf(i18n.T("sysctl.config"), strings.TrimRight(stdout, "\n")))

	vmSwappiness, err := osinfoutil.CheckVMSwappiness(log.Sugar)
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("sysctl.check_fail"), "vm.swappiness", err.Error()))
		return
	}
	console.OK(fmt.Sprintf(i18n.T("sysctl.param_value"), "vm.swappiness", vmSwappiness))

	vmMaxMapCount, err := osinfoutil.CheckVMMaxMapCount(log.Sugar)
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("sysctl.check_fail"), "vm.max_map_count", err.Error()))
		return
	}
	console.OK(fmt.Sprintf(i18n.T("sysctl.param_value"), "vm.max_map_count", vmMaxMapCount))

	netIPLocalPortRange, err := osinfoutil.CheckNetIPLocalPortRange(log.Sugar)
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("sysctl.check_fail"), "net.ipv4.ip_local_port_range", err.Error()))
		return
	}
	console.OK(fmt.Sprintf(i18n.T("sysctl.param_value"), "net.ipv4.ip_local_port_range", netIPLocalPortRange))

	kernelCorePattern, err := osinfoutil.CheckKernelCorePattern(log.Sugar)
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("sysctl.check_fail"), "kernel.core_pattern", err.Error()))
		return
	}
	console.OK(fmt.Sprintf(i18n.T("sysctl.param_value"), "kernel.core_pattern", kernelCorePattern))

	console.Warn(i18n.T("sysctl.check_warn"))
}
