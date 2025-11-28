package checkos

import (
	"fmt"
	"runtime"

	"preinstall/defines/confdef"
	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/check"
	"preinstall/internal/modules/modulecommons/console"

	"github.com/shirou/gopsutil/v3/cpu"
)

func CheckCPU() {
	console.Check(i18n.T("cpu.check"))
	c, err := cpu.Info()
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("cpu.get_info_fail"), err))
		return
	}
	if len(c) == 0 {
		console.Fail(i18n.T("cpu.no_info"))
		return
	}
	info := c[0]
	info.Cores = int32(runtime.NumCPU())
	console.OK(fmt.Sprintf(i18n.T("cpu.vendor"), info.VendorID))
	console.OK(fmt.Sprintf(i18n.T("cpu.model"), info.ModelName))
	console.OK(fmt.Sprintf(i18n.T("cpu.cores"), info.Cores))

	minCores := int32(confdef.Conf().Limit.Hardware.CPUMinCores)
	if info.Cores < minCores {
		console.Warn(fmt.Sprintf(i18n.T("cpu.cores_warn"), info.Cores, minCores))
		check.AddCheckCount()
	}
}
