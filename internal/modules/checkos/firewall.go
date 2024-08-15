package checkos

import (
	"fmt"

	"preinstall/defines/runtimedef"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sysutil/osinfoutil"
	"preinstall/utils/sysutil/userutil"
)

func CheckFirewall() {
	console.Check("防火墙配置")
	active, status, err := osinfoutil.GetFirewallStatus(log.Sugar, runtimedef.GetOSPlatform())
	if err != nil {
		console.Fail(fmt.Sprintf("获取防火墙状态失败: %s", err))
	}
	console.OK(fmt.Sprintf("防火墙状态[systemctl is-active %s]: %s", osinfoutil.GetFirewallName(runtimedef.GetOSPlatform()), status))

	if active {
		defer console.Warn("防火墙开启，相关配置可能会影响网络通信，建议配置合适的端口规则或者关闭防火墙")
		if !userutil.IsCurrentUserRoot() {
			console.Fail("权限不足，无法查看防火墙规则，请使用root用户或者sudo执行")
			return
		}
		output, err := osinfoutil.ShowIPTables(log.Sugar)
		if err != nil {
			console.Fail(fmt.Sprintf("获取防火墙规则失败: %s", err))
		}
		console.OK(fmt.Sprintf("防火墙规则[iptables -L]: \n%s", output))
	}
}
