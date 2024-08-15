package osinfoutil

import (
	"strings"

	"preinstall/defines/bashdef"
	"preinstall/utils/executil"

	"git.yasdb.com/go/yaslog"
)

const (
	STATUS_ACTIVE = "active"
)

var (
	_iptables = bashdef.CMD_IPTABLES
)

func SetIPTablesCommand(command string) {
	_iptables = command
}

func GetFirewallName(platform string) string {
	switch strings.ToLower(platform) {
	case "ubuntu":
		return "ufw"
	default:
		return "firewalld"
	}
}

func GetFirewallStatus(log yaslog.YasLog, platform string) (bool, string, error) {
	name := GetFirewallName(platform)
	return GetServiceStatus(log, name)
}

func ShowIPTables(log yaslog.YasLog) (string, error) {
	execer := executil.NewExecer(log)
	ret, stdout, stderr := execer.Exec(_iptables, "-L")
	if ret != 0 {
		return "", executil.GenerateError(stdout, stderr)
	}
	return stdout, nil
}
