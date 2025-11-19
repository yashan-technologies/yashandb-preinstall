package checkos

import (
	"fmt"

	"preinstall/defines/confdef"
	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/check"
	"preinstall/internal/modules/modulecommons/console"

	"git.yasdb.com/go/yasutil/size"
	"github.com/shirou/gopsutil/v3/mem"
)

func CheckMemory() {
	console.Check(i18n.T("memory.check"))
	v, err := mem.VirtualMemory()
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("memory.get_info_fail"), err))
		return
	}

	console.OK(fmt.Sprintf(i18n.T("memory.total"), size.GenHumanReadableSize(float64(v.Total), 2)))
	console.OK(fmt.Sprintf(i18n.T("memory.available"), size.GenHumanReadableSize(float64(v.Available), 2)))

	minMemory := uint64(confdef.Conf().Limit.Hardware.MemoryMinGB) * 1024 * 1024 * 1024
	if v.Total < minMemory {
		console.Warn(fmt.Sprintf(i18n.T("memory.warn"),
			size.GenHumanReadableSize(float64(v.Total), 2), fmt.Sprintf("%dG", confdef.Conf().Limit.Hardware.MemoryMinGB)))
		check.AddCheckCount()
	}
}
