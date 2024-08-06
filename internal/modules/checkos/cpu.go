package checkos

import (
	"fmt"
	"runtime"

	"preinstall/defines/confdef"
	"preinstall/internal/modules/modulecommons/check"
	"preinstall/internal/modules/modulecommons/console"

	"github.com/shirou/gopsutil/v3/cpu"
)

func CheckCPU() {
	console.Check("CPU处理器")
	c, err := cpu.Info()
	if err != nil {
		console.Fail(fmt.Sprintf("获取CPU信息失败: %s", err))
		return
	}
	if len(c) == 0 {
		console.Fail("未获取到CPU信息")
		return
	}
	info := c[0]
	info.Cores = int32(runtime.NumCPU())
	console.OK("厂商：" + info.VendorID)
	console.OK("名称：" + info.ModelName)
	console.OK("核心：" + fmt.Sprint(info.Cores) + " Cores")

	minCores := int32(confdef.Conf().Limit.Hardware.CPUMinCores)
	if info.Cores < minCores {
		console.Warn(fmt.Sprintf("设备CPU核心数为：%d，建议您使用 %s 核心及以上配置的机器", info.Cores, console.Green.Sprint(minCores)))
		check.AddCheckCount()
	}
}
