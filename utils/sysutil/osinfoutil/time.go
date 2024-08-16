package osinfoutil

import (
	"strings"
	"time"

	"preinstall/utils/executil"

	"git.yasdb.com/go/yaslog"
)

const (
	TZ_SHANGHAI = "Asia/Shanghai"
)

var (
	_timedatectl = "timedatectl"
)

func SetTimedatectlCommand(cmd string) {
	_timedatectl = cmd
}

func GetTimezone(log yaslog.YasLog) (string, error) {
	exexer := executil.NewExecer(log)
	ret, stdout, stderr := exexer.Exec(_timedatectl, "status", "--no-pager")
	if ret != 0 {
		return "", executil.GenerateError(stdout, stderr)
	}
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Time zone") {
			return strings.TrimSpace(strings.Split(line, ":")[1]), nil
		}
	}
	return stdout, nil
}

func IsTimeUTC8() bool {
	_, offset := time.Now().Zone()
	return offset == 8*60*60
}

func SetTimezone(log yaslog.YasLog, timezone string) error {
	exexer := executil.NewExecer(log)
	ret, stdout, stderr := exexer.Exec(_timedatectl, "set-timezone", timezone)
	if ret != 0 {
		return executil.GenerateError(stdout, stderr)
	}
	return nil
}

func SetNtp(log yaslog.YasLog, enable bool) error {
	exexer := executil.NewExecer(log)
	var value string
	if enable {
		value = "true"
	} else {
		value = "false"
	}
	ret, stdout, stderr := exexer.Exec(_timedatectl, "set-ntp", value)
	if ret != 0 {
		return executil.GenerateError(stdout, stderr)
	}
	return nil
}
