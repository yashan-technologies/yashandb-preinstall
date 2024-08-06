package iniutil

import (
	"fmt"

	"preinstall/commons/consts"

	"git.yasdb.com/go/yasutil/fs"
	"gopkg.in/ini.v1"
)

var (
	_loadOptions = ini.LoadOptions{
		SkipUnrecognizableLines: true,  // 忽略无法识别的行
		IgnoreInlineComment:     false, // 保留行内注释
	}
)

func LoadIni(fname string) (*ini.File, error) {
	if !fs.IsFileExist(fname) {
		return nil, fmt.Errorf("file not exist: %s", fname)
	}
	return ini.LoadSources(_loadOptions, fname)
}

func WriteIni(fname string, data map[string]map[string]string /*map[section]map[key]value*/, backup bool) (err error) {
	if len(data) == 0 {
		return nil
	}

	if backup {
		backupFile := fname + consts.BakupExt()
		if err = fs.CopyFile(fname, backupFile); err != nil {
			return
		}
	}

	var iniFile *ini.File
	if !fs.IsFileExist(fname) {
		iniFile = ini.Empty()
	} else {
		iniFile, err = ini.LoadSources(_loadOptions, fname)
		if err != nil {
			return
		}
	}

	for section, kvs := range data {
		if iniFile.HasSection(section) {
			for key, value := range kvs {
				iniFile.Section(section).Key(key).SetValue(value)
			}
		} else {
			se, err := iniFile.NewSection(section)
			if err != nil {
				return err
			}
			for key, value := range kvs {
				se.NewKey(key, value)
			}
		}
	}
	return iniFile.SaveTo(fname)
}
