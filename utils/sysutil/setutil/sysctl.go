package setutil

import (
	"strings"

	"preinstall/utils/executil"
	"preinstall/utils/iniutil"
	"preinstall/utils/sysutil/fileutil"

	"git.yasdb.com/go/yaslog"
)

const (
	sysctl_conf_path = "/etc/sysctl.conf"
)

const (
	hash  = "#"
	equal = "="
)

var _sysctl = "sysctl"

func SetSysctlCommand(sysctl string) {
	_sysctl = sysctl
}

// /etc/sysctl.conf和/etc/sysctl.d/下的文件都会生效
// 但是/etc/sysctl.conf的优先级更高，并且后面的配置会覆盖前面的配置
// 例如：/etc/sysctl.conf中配置了vm.swappiness=10 vm.swappiness=15，/etc/sysctl.d/a.conf中配置了vm.swappiness=20
// 那么最终生效的是vm.swappiness=15，通过sysctl --system会加载所有的配置并打印，再执行sysctl vm.swappiness可以查看最终的值
// sysctl -p会加载/etc/sysctl.conf，sysctl --system 会加载/etc/sysctl.conf和/etc/sysctl.d/下的所有文件
func UpdateSysctlConf(log yaslog.YasLog, values ...string) error {
	kd, err := iniutil.LoadKeyData(sysctl_conf_path, "")
	if err != nil {
		return err
	}

	data := make(map[string]string)
	for _, value := range values {
		items := strings.Split(value, equal)
		key := strings.TrimSpace(items[0])
		value := strings.TrimSpace(strings.TrimPrefix(value, items[0]+equal))
		data[key] = value
	}

	var lines []string
	for _, line := range kd.RawLines {
		if len(line) == 0 || strings.HasPrefix(line, hash) {
			lines = append(lines, line)
			continue
		}

		// 删除已存在的配置
		if strings.Contains(line, equal) {
			key := strings.TrimSpace(strings.Split(line, equal)[0])
			if _, ok := data[key]; ok {
				continue
			}
		}
		lines = append(lines, line)
	}

	content := strings.Join(lines, newline)
	content = strings.TrimRight(content, newline)
	content += newline + newline
	content += strings.Join(values, newline) + newline

	if err := fileutil.BackupFile(sysctl_conf_path); err != nil {
		return err
	}

	if err := fileutil.WriteFile(sysctl_conf_path, []byte(content)); err != nil {
		return err
	}

	_, err = ReadSysctlConf(log)
	return err
}

func ReadSysctlConf(log yaslog.YasLog) (string, error) {
	execer := executil.NewExecer(log)
	ret, stdout, stderr := execer.Exec(_sysctl, "-p")
	if ret != 0 {
		return "", executil.GenerateError(stdout, stderr)
	}
	return stdout, nil
}
