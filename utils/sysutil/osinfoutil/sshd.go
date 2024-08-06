package osinfoutil

import (
	"strings"

	"git.yasdb.com/go/yaslog"
)

func GetSshServiceName(platform string) string {
	switch strings.ToLower(platform) {
	case "ubuntu":
		return "ssh"
	default:
		return "sshd"
	}
}

func GetSshdStatus(log yaslog.YasLog, platform string) (bool, string, error) {
	name := GetSshServiceName(platform)
	return GetServiceStatus(log, name)
}
