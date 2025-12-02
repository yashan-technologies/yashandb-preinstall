package main

import (
	"fmt"
	"strings"

	"preinstall/commons/flags"
	"preinstall/commons/std"
	"preinstall/defines/compiledef"
	"preinstall/defines/confdef"
	"preinstall/defines/runtimedef"
	"preinstall/i18n"
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
	if err := initApp(&app); err != nil {
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

func initApp(app *internal.App) error {
	// 初始化运行时环境
	if err := runtimedef.InitRuntime(); err != nil {
		return err
	}

	// 加载配置文件
	if err := confdef.InitConfig(app.Config); err != nil {
		return err
	}

	// 确定最终使用的语言：命令行参数优先，否则使用配置文件中的设置，最后默认为 "zh"
	language := app.Language
	if language == "" {
		// 命令行未指定，使用配置文件中的设置
		language = confdef.Conf().Language
	}
	if language == "" {
		// 配置文件也未设置，使用默认值
		language = "zh"
	}
	// 初始化i18n
	if err := i18n.Init(language); err != nil {
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
