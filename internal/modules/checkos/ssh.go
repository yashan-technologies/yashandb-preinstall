package checkos

import (
	"fmt"

	"preinstall/defines/confdef"
	"preinstall/defines/runtimedef"
	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sshutil"
	"preinstall/utils/sysutil/osinfoutil"
)

func CheckSSH() {
	console.Check(i18n.T("ssh.check"))
	active, status, err := osinfoutil.GetSshdStatus(log.Sugar, runtimedef.GetOSPlatform())
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("ssh.get_status_fail"), err))
	}
	console.OK(fmt.Sprintf(i18n.T("ssh.status"), status))

	if !active {
		console.Warn(i18n.T("ssh.service_warn"))
	}

	if len(confdef.YashanDBConf().Hosts) > 0 {
		console.Info(i18n.T("ssh.testing"))
	}
	for _, h := range confdef.YashanDBConf().Hosts {
		ssher := sshutil.NewSSH(h.IP, h.Port, h.User, h.Password)
		if err := ssher.CheckSSHConnection(); err != nil {
			console.Fail(fmt.Sprintf(i18n.T("ssh.connect_fail"), h.IP, err))
			continue
		}
		console.OK(fmt.Sprintf(i18n.T("ssh.connect_success"), h.IP))
	}
}
