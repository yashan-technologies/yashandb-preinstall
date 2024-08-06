package setutil

import (
	"fmt"
	"os/exec"

	"preinstall/utils/executil"
	"preinstall/utils/sysutil/fileutil"

	"git.yasdb.com/go/yaslog"
)

const (
	grub_cfg  = "grub.cfg"
	base_path = "/boot"
)

var (
	_updateGrub   = "update-grub"
	_grubMkConfig = []string{"grub2-mkconfig", "grub-mkconfig"}
	_grubCfg      string
)

var ErrMultipleGrubCfg = fmt.Errorf("multiple grub.cfg found")

func SetUpdateGrubCommand(command string) {
	_updateGrub = command
}

func SetGrubMkConfigCommand(commands []string) {
	_grubMkConfig = commands
}

func SetGrubCfg(file string) {
	_grubCfg = file
}

func UpdateGrubConfig(log yaslog.YasLog) ([]string, error) {
	execer := executil.NewExecer(log)
	updateGrubExists := isCommandExists(_updateGrub)
	if updateGrubExists {
		ret, stdout, stderr := execer.Exec(_updateGrub)
		if ret != 0 {
			return nil, executil.GenerateError(stdout, stderr)
		}
		return nil, nil
	}

	grubMkConfig := findGrubMkConfigCommand()
	if len(grubMkConfig) == 0 {
		return nil, fmt.Errorf("%s command not found", _grubMkConfig)
	}

	grubConf := _grubCfg
	if len(grubConf) == 0 {
		files, err := fileutil.FindFile(base_path, grub_cfg)
		if err != nil {
			return nil, err
		}
		if len(files) == 0 {
			return nil, fmt.Errorf("grub.cfg not found in %s", base_path)
		}
		if len(files) > 1 {
			return files, ErrMultipleGrubCfg
		}
		grubConf = files[0]
	}
	ret, stdout, stderr := execer.Exec(grubMkConfig, "-o", grubConf)
	if ret != 0 {
		return []string{grubConf}, executil.GenerateError(stdout, stderr)
	}
	return []string{grubConf}, nil
}

func findGrubMkConfigCommand() string {
	for _, command := range _grubMkConfig {
		if isCommandExists(command) {
			return command
		}
	}
	return ""
}

func isCommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
