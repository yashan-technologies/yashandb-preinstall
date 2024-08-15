package setutil

import (
	"fmt"
	"os"
	"strings"

	"preinstall/defines/bashdef"
	"preinstall/utils/executil"

	"git.yasdb.com/go/yaslog"
)

const (
	ENABLE_STATUS_ALWAYS  = "always"
	ENABLE_STATUS_MADVISE = "madvise"
	ENABLE_STATUS_NEVER   = "never"
)

const transparent_hugepage_path = "/sys/kernel/mm/transparent_hugepage/"

var (
	_echo = bashdef.CMD_ECHO
	_bash = bashdef.CMD_BASH
)

func SetEchoCommand(echo string) {
	_echo = echo
}

func SetBashCommand(bash string) {
	_bash = bash
}

func GetTransparentHugePageSetting(setting string) (string, error) {
	content, err := os.ReadFile(transparent_hugepage_path + setting)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// 直接修改/sys/kernel/mm/transparent_hugepage/enabled文件通常不会生效，
// 因为该文件是内核接口的一部分，它是用来通过特定的接口与内核交互的，而不是像普通文件那样用于存储数据。
// 因此，要改变透明大页的设置，你需要使用echo命令（或其他类似的工具，如printf）来将特定的字符串（例如[always], [never], [madvise]）写入到这个文件中。
// 这样，内核就会接收到信号，并根据写入的内容改变透明大页的行为。
func SetTransparentHugePageSetting(log yaslog.YasLog, setting, value string) error {
	originFile := transparent_hugepage_path + setting
	execer := executil.NewExecer(log)
	cmd := fmt.Sprintf("%s %s > %s", _echo, value, originFile)
	ret, stdout, stderr := execer.Exec(_bash, "-c", cmd)
	if ret != 0 {
		return executil.GenerateError(stdout, stderr)
	}
	return nil
}

// 设置透明大页的设置，临时生效重启后失效
func SetTransparentHugePageEnabledToNever(log yaslog.YasLog) error {
	return SetTransparentHugePageSetting(log, "enabled", ENABLE_STATUS_NEVER)
}

func GetTransparentHugePageEnabled() (string, error) {
	content, err := GetTransparentHugePageSetting("enabled")
	if err != nil {
		return "", err
	}
	// 解析内容，根据空格分隔和[]匹配
	// 例如：[always] madvise never，解析出always
	return parseTransparentHugePageSetting(content), nil
}

func IsTransparentHugePageEnabledNever() (bool, error) {
	enabled, err := GetTransparentHugePageEnabled()
	if err != nil {
		return false, err
	}
	return enabled == ENABLE_STATUS_NEVER, nil
}

func parseTransparentHugePageSetting(content string) string {
	parts := strings.Fields(content)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]") {
			return strings.TrimPrefix(strings.TrimSuffix(part, "]"), "[")
		}
	}
	return ""
}
