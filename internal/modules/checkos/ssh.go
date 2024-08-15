package checkos

import (
	"fmt"

	"preinstall/defines/confdef"
	"preinstall/defines/runtimedef"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sshutil"
	"preinstall/utils/sysutil/osinfoutil"
)

func CheckSSH() {
	console.Check("SSH服务")
	active, status, err := osinfoutil.GetSshdStatus(log.Sugar, runtimedef.GetOSPlatform())
	if err != nil {
		console.Fail(fmt.Sprintf("获取SSH服务状态失败: %s", err))
	}
	console.OK(fmt.Sprintf("SSH服务状态[systemctl is-active %s]: %s", osinfoutil.GetSshServiceName(runtimedef.GetOSPlatform()), status))

	if !active {
		console.Warn("SSH服务异常，请检查")
	}

	if len(confdef.YashanDBConf().Hosts) > 0 {
		console.Info("SSH连接测试......")
	}
	for _, h := range confdef.YashanDBConf().Hosts {
		ssher := sshutil.NewSSH(h.IP, h.Port, h.User, h.Password)
		if err := ssher.CheckSSHConnection(); err != nil {
			console.Fail(fmt.Sprintf("SSH连接[%s]失败: %s", h.IP, err))
		}
		console.OK(fmt.Sprintf("SSH连接[%s]成功", h.IP))
	}
}
