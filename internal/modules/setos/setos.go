package setos

import (
	"fmt"

	"preinstall/internal/modules/modulecommons/console"
	"preinstall/internal/modules/setos/setconf"
	"preinstall/internal/modules/setos/setuser"
	"preinstall/internal/modules/setos/setyasdb"
)

func Set(setDiskQueneScheduler bool) {
	fmt.Println(console.White.Sprint("<========================开始配置系统========================>\n"))
	setuser.SetUser()
	fmt.Println()

	setconf.DisableTransparentHugePageAndNUMA()
	fmt.Println()

	setconf.SetUserLimits()
	fmt.Println()

	setconf.UpdateSysctl()
	fmt.Println()

	setyasdb.SetYashanDBPath(setDiskQueneScheduler)
	fmt.Println()
}
