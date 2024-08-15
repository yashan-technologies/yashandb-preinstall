package checkos

import (
	"strings"

	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sysutil/osinfoutil"
	"preinstall/utils/sysutil/setutil"
)

func CheckSysctlConf() {
	console.Check("检查内核参数配置")

	stdout, err := setutil.ReadSysctlConf(log.Sugar)
	if err != nil {
		console.Fail("读取内核参数配置失败：" + err.Error())
		return
	}

	console.OK("内核参数配置[sysctl -p]：\n" + strings.TrimRight(stdout, "\n"))

	vmSwappiness, err := osinfoutil.CheckVMSwappiness(log.Sugar)
	if err != nil {
		console.Fail("检查vm.swappiness配置失败：" + err.Error())
		return
	}
	console.OK("[cat /proc/sys/vm/swappiness]：" + vmSwappiness)

	vmMaxMapCount, err := osinfoutil.CheckVMMaxMapCount(log.Sugar)
	if err != nil {
		console.Fail("检查vm.max_map_count配置失败：" + err.Error())
		return
	}
	console.OK("[cat /proc/sys/vm/max_map_count]：" + vmMaxMapCount)

	netIPLocalPortRange, err := osinfoutil.CheckNetIPLocalPortRange(log.Sugar)
	if err != nil {
		console.Fail("检查net.ipv4.ip_local_port_range配置失败：" + err.Error())
		return
	}
	console.OK("[cat /proc/sys/net/ipv4/ip_local_port_range]：" + netIPLocalPortRange)

	kernelCorePattern, err := osinfoutil.CheckKernelCorePattern(log.Sugar)
	if err != nil {
		console.Fail("检查kernel.core_pattern配置失败：" + err.Error())
		return
	}
	console.OK("[cat /proc/sys/kernel/core_pattern]：" + kernelCorePattern)

	console.Warn("请根据实际情况检查内核参数配置是否符合要求，特别是CoreDump配置（kernel.core_pattern）\n需要确保路径存在并且权限正确，以及考虑是否影响其他用户的应用程序")
}
