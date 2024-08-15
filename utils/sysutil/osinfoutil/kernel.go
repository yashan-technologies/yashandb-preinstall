package osinfoutil

import (
	"strings"

	"preinstall/utils/executil"

	"git.yasdb.com/go/yaslog"
)

var (
	_cat = "cat"
)

func SetCatCommand(cat string) {
	_cat = cat
}

func CheckVMSwappiness(log yaslog.YasLog) (string, error) {
	return catFile(log, "/proc/sys/vm/swappiness")
}

func CheckVMMaxMapCount(log yaslog.YasLog) (string, error) {
	return catFile(log, "/proc/sys/vm/max_map_count")
}

func CheckNetIPLocalPortRange(log yaslog.YasLog) (string, error) {
	return catFile(log, "/proc/sys/net/ipv4/ip_local_port_range")
}

func CheckKernelCorePattern(log yaslog.YasLog) (string, error) {
	return catFile(log, "/proc/sys/kernel/core_pattern")
}

func catFile(log yaslog.YasLog, file string) (string, error) {
	execer := executil.NewExecer(log)
	ret, stdout, stderr := execer.Exec(_cat, file)
	if ret != 0 {
		return "", executil.GenerateError(stdout, stderr)
	}
	return strings.TrimRight(stdout, "\n"), nil
}
