package internal

import (
	"fmt"
	"path"

	"preinstall/commons/flags"
	"preinstall/commons/std"
	"preinstall/defines/runtimedef"
	"preinstall/internal/modules/checkos"
	"preinstall/internal/modules/setos"
)

type App struct {
	flags.Globals
	CheckOnly             bool `name:"check-only"               short:"o" help:"跳过配置操作，只进行环境检查."`
	EnableIOTest          bool `name:"io"                       short:"i" help:"测试IO性能."`
	SetDiskQueneScheduler bool `name:"set-disk-quene-scheduler" short:"s" help:"设置磁盘队列调度器."`
}

func (a *App) Preinstall() {
	if !a.CheckOnly {
		setos.Set(a.SetDiskQueneScheduler)
	}
	checkos.Check(a.EnableIOTest)
	fmt.Printf("部署前配置检查已完成，请检查失败或告警信息（结果存放路径：%s.*）。\n",
		path.Join(runtimedef.ResultPath(), std.CONSOLE_OUT_FILE))
	if !a.EnableIOTest {
		fmt.Println("另外，如果有需要您可以使用 -i 参数在安装路径磁盘进行I/O性能测试。")
	}
}
