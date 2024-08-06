package runtimedef

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

const (
	_ENV_HOME = "YASHANDB_PREINSTALL_HOME"
)

const (
	_DIR_CONF    = "config"
	_DIR_LOG     = "log"
	_DIR_PLUGINS = "plugins"
)

var _home string

func HomePath() string {
	return _home
}

func ConfPath() string {
	return path.Join(_home, _DIR_CONF)
}

func LogPath() string {
	return path.Join(_home, _DIR_LOG)
}

func PluginsPath() string {
	return path.Join(_home, _DIR_PLUGINS)
}

func Fio() string {
	return path.Join(PluginsPath(), "fio")
}

func setHome(v string) {
	_home = v
}

func getHomeFromEnv() string {
	return os.Getenv(_ENV_HOME)
}

func genHomeFromEnv() (home string, err error) {
	homeFromEnv := getHomeFromEnv()
	if len(homeFromEnv) > 0 {
		homeFromEnv, err = filepath.Abs(homeFromEnv)
		if err != nil {
			return
		}
		home = homeFromEnv
		return
	}
	return
}

func genHomeFromRelativePath() (home string, err error) {
	executeable, err := getExecutable()
	if err != nil {
		return
	}
	home, err = filepath.Abs(path.Dir(path.Dir(executeable)))
	return
}

func initHomePath() error {
	// 首先尝试从环境变量中获取
	home, err := genHomeFromEnv()
	if err != nil {
		fmt.Printf("[WARN] get home from env: %v\n", err)
	} else if len(home) > 0 {
		setHome(home)
		return nil
	}

	// 如果环境变量中没有，尝试从相对路径中获取
	home, err = genHomeFromRelativePath()
	if err != nil {
		return err
	}
	setHome(home)
	return nil
}
