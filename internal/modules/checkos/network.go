package checkos

import (
	"fmt"
	"strings"

	"preinstall/defines/confdef"
	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/utils/sysutil/osinfoutil"
)

func CheckNetwork() {
	console.Check(i18n.T("network.check"))
	interfaces, err := osinfoutil.GetPhysicalNetInterfaces()
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("network.get_info_fail"), err))
		return
	}
	for _, i := range interfaces {
		var addrs []string
		for _, addr := range i.Addrs {
			addrs = append(addrs, addr.String())
		}
		var speed string
		if i.Speed <= 0 {
			speed = "Unknown"
		} else {
			speed = fmt.Sprintf("%d Mbps", i.Speed)
		}
		console.OK(fmt.Sprintf(i18n.T("network.interface_info"), i.Name, strings.Join(addrs, ", "), speed))

		if i.Speed <= 0 {
			console.Warn(fmt.Sprintf(i18n.T("network.speed_unknown"), i.Name))
		} else if i.Speed < confdef.Conf().Limit.Hardware.NetworkMinBandWidthMbps {
			console.Warn(fmt.Sprintf(i18n.T("network.speed_warn"), fmt.Sprintf("%dMbps", confdef.Conf().Limit.Hardware.NetworkMinBandWidthMbps)))
		}
	}
}
