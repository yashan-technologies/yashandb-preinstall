package checkos

import (
	"fmt"

	"preinstall/defines/runtimedef"
	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sysutil/osinfoutil"
	"preinstall/utils/sysutil/userutil"
)

func CheckFirewall() {
	console.Check(i18n.T("firewall.check"))
	active, status, err := osinfoutil.GetFirewallStatus(log.Sugar, runtimedef.GetOSPlatform())
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("firewall.get_status_fail"), err))
	}
	console.OK(fmt.Sprintf(i18n.T("firewall.status"), status))

	if active {
		defer console.Warn(i18n.T("firewall.warn"))
		if !userutil.IsCurrentUserRoot() {
			console.Fail(i18n.T("firewall.no_privilege"))
			return
		}
		output, err := osinfoutil.ShowIPTables(log.Sugar)
		if err != nil {
			console.Fail(fmt.Sprintf(i18n.T("firewall.get_rules_fail"), err))
		}
		console.OK(fmt.Sprintf(i18n.T("firewall.rules"), output))
	}
}
