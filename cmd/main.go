package main

import (
	"fmt"
	"strings"

	"preinstall/commons/flags"
	"preinstall/commons/std"
	"preinstall/defines/compiledef"
	"preinstall/defines/confdef"
	"preinstall/defines/runtimedef"
	"preinstall/internal"
	"preinstall/log"

	"github.com/alecthomas/kong"
)

const (
	_APP_NAME        = "preinstall"
	_APP_DESCRIPTION = "YashanDB Install Preparation."
)

func main() {
	var app internal.App
	options := flags.NewAppOptions(_APP_NAME, _APP_DESCRIPTION, compiledef.GetAPPVersion())
	ctx := kong.Parse(&app, options...)
	if err := initApp(app); err != nil {
		ctx.FatalIfErrorf(err)
	}
	finalize := std.GetRedirecter().RedirectStd()
	defer finalize()
	std.WriteToFile(fmt.Sprintf("$ %s %s\n", _APP_NAME, strings.Join(ctx.Args, " ")))
	app.Preinstall()
}

func initLogger(logPath, level string) error {
	optFuncs := []log.OptFunc{
		log.SetLogPath(logPath),
		log.SetLevel(level),
	}
	return log.InitLogger(_APP_NAME, log.NewLogOption(optFuncs...))
}

func initApp(app internal.App) error {
	// 首先初始化运行时环境
	if err := runtimedef.InitRuntime(); err != nil {
		return err
	}

	// 再加载配置文件
	if err := confdef.InitConfig(app.Config); err != nil {
		return err
	}

	// 根据运行时环境和配置文件初始化日志
	if err := initLogger(runtimedef.LogPath(), confdef.Conf().LogLevel); err != nil {
		return err
	}
	log.Sugar.Debugf("Conf: %s", confdef.Conf().ToJSON())
	log.Sugar.Debugf("YashanDB Conf: %s", confdef.YashanDBConf().ToJSON())

	// 初始化终端重定向器
	if err := std.InitRedirecter(); err != nil {
		return err
	}
	return nil
}
