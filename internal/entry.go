package internal

import (
	"fmt"
	"path"

	"preinstall/commons/flags"
	"preinstall/commons/std"
	"preinstall/defines/runtimedef"
	"preinstall/i18n"
	"preinstall/internal/modules/checkos"
	"preinstall/internal/modules/setos"
)

type App struct {
	flags.Globals
	CheckOnly             bool `name:"check-only"               short:"o" help:"Skip configuration operations, only perform environment checks."`
	EnableIOTest          bool `name:"io"                       short:"i" help:"Test IO performance."`
	SetDiskQueneScheduler bool `name:"set-disk-quene-scheduler" short:"s" help:"Set disk queue scheduler."`
}

func (a *App) Preinstall() {
	if !a.CheckOnly {
		setos.Set(a.SetDiskQueneScheduler)
	}
	checkos.Check(a.EnableIOTest)
	fmt.Printf(i18n.T("main.complete")+"\n",
		path.Join(runtimedef.ResultPath(), std.CONSOLE_OUT_FILE))
	if !a.EnableIOTest {
		fmt.Println(i18n.T("main.io_tip"))
	}
}
