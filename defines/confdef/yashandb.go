package confdef

import (
	"encoding/json"
	"path"

	"preinstall/defines/errdef"
	"preinstall/defines/runtimedef"

	"git.yasdb.com/go/yasutil/fs"
	"github.com/BurntSushi/toml"
)

const yashandb_conf = "yashandb.toml"

var _yashanDB YashanDB

type YashanDB struct {
	InstallPath      string   `toml:"install_path"`
	YasdbHome        string   `toml:"yasdb_home"`
	YasdbData        string   `toml:"yasdb_data"`
	YasdbBack        string   `toml:"yasdb_back"`
	YasdbBackSubdirs []string `toml:"yasdb_back_subdirs"`
	Hosts            []Host   `toml:"hosts"`
}

type Host struct {
	IP       string `toml:"ip"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Port     int    `toml:"port"`
}

func YashanDBConf() YashanDB {
	return _yashanDB
}

func (c YashanDB) ToJSON() string {
	data, _ := json.MarshalIndent(c, "", "    ")
	return string(data)
}

func initYashanDBConfig(confPath string) error {
	if !path.IsAbs(confPath) {
		confPath = path.Join(runtimedef.ConfPath(), confPath)
	}
	if !fs.IsFileExist(confPath) {
		return &errdef.ErrFileNotFound{FName: confPath}
	}
	if _, err := toml.DecodeFile(confPath, &_yashanDB); err != nil {
		return &errdef.ErrFileParseFailed{FName: confPath, Err: err}
	}
	return nil
}
