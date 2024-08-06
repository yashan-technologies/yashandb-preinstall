package setutil

import (
	"os"
	"path"
	"regexp"
	"strings"

	"preinstall/commons/consts"

	"git.yasdb.com/go/yasutil/fs"
)

const (
	serurity_limits_conf     = "/etc/security/limits.conf"
	serurity_limits_conf_dir = "/etc/security/limits.d"
	comment_prefix           = "# Limits of YashanDB User: "
	newline                  = "\n"
)

// /etc/security/limits.d文件下也可以配置用户限制
// 并且优先级会比/etc/security/limits.conf高
// 在/etc/security/limits.conf中配置后，也会被/etc/security/limits.d下的文件覆盖
// 并且如果在/etc/security/limits.d下有多个文件同时配置了同一个用户的限制，不确定哪个会生效
// 例如：现在a.conf中限制了 user1 limit1 value1，b.conf中限制了 user1 limit1 value2
// 那么最终生效的是哪个值，不确定
// 所以需要删除所有配置文件中的用户限制，然后重新配置
func SetUserLimits(username string, limits []string, backup bool) error {
	// 删除已存在的限制
	if err := deleteUserLimits(username, backup); err != nil {
		return err
	}

	var lines []string
	limitFile := path.Join(serurity_limits_conf_dir, username+".conf")
	exists := fs.IsFileExist(limitFile)
	if exists {
		content, err := os.ReadFile(limitFile)
		if err != nil {
			return err
		}
		lines = strings.Split(string(content), newline)
	}
	lines = append(lines, comment_prefix+username)
	lines = append(lines, limits...)
	lines = append(lines, newline)
	content := strings.Join(lines, newline)
	content = newline + content

	if exists && backup {
		backupFile := limitFile + consts.BakupExt()
		if err := fs.CopyFile(limitFile, backupFile); err != nil {
			return err
		}
	}
	// 正则替换content，把连续多个\n替换成一个\n
	content = regexp.MustCompile(`\n+`).ReplaceAllString(content, newline)
	// 去除开始的\n
	content = strings.TrimLeft(content, newline)
	// 如果文件不是\n结尾，加上\n
	if !strings.HasSuffix(content, newline) {
		content += newline
	}
	return os.WriteFile(limitFile, []byte(content), 0644)
}

func CheckUserLimits() {

}

func deleteUserLimits(username string, backup bool) error {
	var confs []string
	confs = append(confs, serurity_limits_conf)

	// 查询/etc/security/limits.d下的文件
	files, err := os.ReadDir(serurity_limits_conf_dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		// 只处理.conf文件
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".conf") {
			continue
		}
		confs = append(confs, path.Join(serurity_limits_conf_dir, file.Name()))
	}

	for _, conf := range confs {
		content, err := os.ReadFile(conf)
		if err != nil {
			return err
		}

		var contentLines []string
		var matched bool
		lines := strings.Split(string(content), newline)
		for _, line := range lines {
			// 删除已存在的限制
			trimLine := strings.TrimSpace(line)
			hasPrefixUsername := strings.HasPrefix(trimLine, username)
			hasPrefixComment := strings.HasPrefix(trimLine, comment_prefix)
			if hasPrefixUsername || hasPrefixComment {
				matched = true
				continue
			}
			contentLines = append(contentLines, line)
		}

		if matched {
			if backup {
				backupFile := conf + consts.BakupExt()
				if err := fs.CopyFile(conf, backupFile); err != nil {
					return err
				}
			}
			if err := os.WriteFile(conf, []byte(strings.Join(contentLines, newline)), 0644); err != nil {
				return err
			}
		}

	}

	return nil
}
