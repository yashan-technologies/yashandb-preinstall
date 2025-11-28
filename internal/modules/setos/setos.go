package setos

import (
	"fmt"

	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/internal/modules/setos/setconf"
	"preinstall/internal/modules/setos/setuser"
	"preinstall/internal/modules/setos/setyasdb"
)

func Set(setDiskQueneScheduler bool) {
	fmt.Println(console.White.Sprint(i18n.T("main.start_set") + "\n"))
	setuser.SetUser()
	fmt.Println()

	setyasdb.SetYashanDBPath(setDiskQueneScheduler)
	fmt.Println()

	setconf.SetUserLimits()
	fmt.Println()

	setconf.DisableTransparentHugePageAndNUMA()
	fmt.Println()

	setconf.UpdateSysctl()
	fmt.Println()

	setconf.SetTime()
	fmt.Println()
}
