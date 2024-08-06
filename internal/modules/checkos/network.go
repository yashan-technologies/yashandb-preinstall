package checkos

import (
	"fmt"
	"strings"

	"preinstall/defines/confdef"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/utils/sysutil/osinfoutil"
)

func CheckNetwork() {
	console.Check("网络设备")
	interfaces, err := osinfoutil.GetPhysicalNetInterfaces()
	if err != nil {
		console.Fail(fmt.Sprintf("获取网络接口信息失败: %s", err))
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
		console.OK(fmt.Sprintf("网络接口：%s\n    └─ IP地址：%s\n    └─ 带宽：%s", i.Name, strings.Join(addrs, ", "), speed))

		if i.Speed <= 0 {
			console.Warn(fmt.Sprintf("未获取到 %s 的带宽信息", i.Name))
		} else if i.Speed < confdef.Conf().Limit.Hardware.NetworkMinBandWidthMbps {
			console.Warn(fmt.Sprintf("建议您使用 %s 及以上的网络设备", console.Green.Sprintf("%dMbps", confdef.Conf().Limit.Hardware.NetworkMinBandWidthMbps)))
		}
	}
}
