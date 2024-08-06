package confdef

import (
	"path"

	"preinstall/defines/errdef"
	"preinstall/defines/runtimedef"
	"preinstall/utils/sysutil/osinfoutil"
	"preinstall/utils/sysutil/setutil"
	"preinstall/utils/sysutil/userutil"

	"git.yasdb.com/go/yasutil/fs"
	"github.com/BurntSushi/toml"
)

var _conf Config

func InitConfig(confPath string) error {
	if !path.IsAbs(confPath) {
		confPath = path.Join(runtimedef.HomePath(), confPath)
	}
	if !fs.IsFileExist(confPath) {
		return &errdef.ErrFileNotFound{FName: confPath}
	}
	if _, err := toml.DecodeFile(confPath, &_conf); err != nil {
		return &errdef.ErrFileParseFailed{FName: confPath, Err: err}
	}
	initCommands()

	if err := initYashanDBConfig(yashandb_conf); err != nil {
		return err
	}
	return nil
}

func initCommands() {
	userutil.SetGroupaddCommand(_conf.Commands.Groupadd)
	userutil.SetUseraddCommand(_conf.Commands.Useradd)
	userutil.SetUsermodCommand(_conf.Commands.Usermod)
	userutil.SetSudoCommand(_conf.Commands.Sudo)

	osinfoutil.SetSystemctlCommand(_conf.Commands.Systemctl)
	osinfoutil.SetIPTablesCommand(_conf.Commands.IPTables)
	osinfoutil.SetSuCommand(_conf.Commands.Su)
	osinfoutil.SetUlimitCommand(_conf.Commands.Ulimit)
	osinfoutil.SetCatCommand(_conf.Commands.Cat)
	osinfoutil.SetTimedatectlCommand(_conf.Commands.Timedatectl)

	setutil.SetEchoCommand(_conf.Commands.Echo)
	setutil.SetBashCommand(_conf.Commands.Bash)
	setutil.SetUpdateGrubCommand(_conf.Commands.UpdateGrub)
	setutil.SetGrubMkConfigCommand(_conf.Commands.GrubMkConfig)
	setutil.SetGrubCfg(_conf.Commands.GrubCgf)
	setutil.SetSysctlCommand(_conf.Commands.Sysctl)
	setutil.SetSwapoffCommand(_conf.Commands.Swapoff)
}
