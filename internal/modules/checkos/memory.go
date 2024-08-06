package checkos

import (
	"fmt"

	"preinstall/defines/confdef"
	"preinstall/internal/modules/modulecommons/check"
	"preinstall/internal/modules/modulecommons/console"

	"git.yasdb.com/go/yasutil/size"
	"github.com/shirou/gopsutil/v3/mem"
)

func CheckMemory() {
	console.Check("Memory内存")
	v, err := mem.VirtualMemory()
	if err != nil {
		console.Fail(fmt.Sprintf("获取内存信息失败: %s", err))
		return
	}

	console.OK("内存容量：" + size.GenHumanReadableSize(float64(v.Total), 2))
	console.OK("可用内存：" + size.GenHumanReadableSize(float64(v.Available), 2))

	minMemory := uint64(confdef.Conf().Limit.Hardware.MemoryMinGB) * 1024 * 1024 * 1024
	if v.Total < minMemory {
		console.Warn(fmt.Sprintf("设备内存总量为：%s，建议您使用 %s 及以上配置的机器",
			size.GenHumanReadableSize(float64(v.Total), 2), console.Green.Sprintf("%dG", confdef.Conf().Limit.Hardware.MemoryMinGB)))
		check.AddCheckCount()
	}
}
