package userutil

import (
	"os"
	"os/user"
	"strconv"
	"strings"

	"preinstall/utils/executil"

	"git.yasdb.com/go/yaslog"
)

func GetUsernameById(id int) (string, error) {
	u, err := user.LookupId(strconv.FormatInt(int64(id), 10))
	if err != nil {
		return "", err
	}
	return u.Username, nil
}

func GetUserByName(username string) (*user.User, error) {
	return user.Lookup(username)
}

func GetCurrentUser() (string, error) {
	return GetUsernameById(os.Getuid())
}

func IsCurrentUserRoot() bool {
	return os.Getuid() == 0
}

func IsUserExists(username string) bool {
	_, err := user.Lookup(username)
	return err == nil
}

func AddUserIfNotExists(log yaslog.YasLog, username, primaryGroup, homePath string, supplementaryGroups ...string) (bool, error) {
	var ret int
	var stdout, stderr string
	var exists bool
	execer := executil.NewExecer(log)
	if IsUserExists(username) {
		exists = true
		// 如果用户存在，则只进行添加附加组操作
		ret, stdout, stderr = execer.Exec(_usermodCommand, "-aG", strings.Join(supplementaryGroups, ","), username)
	} else {
		// 如果用户不存在，则添加用户，并设置主组、附加组、家目录
		ret, stdout, stderr = execer.Exec(_useraddCommand, "-m", "-g", primaryGroup, "-G", strings.Join(supplementaryGroups, ","), "-d", homePath, username)
	}
	if ret != 0 {
		return exists, executil.GenerateError(stdout, stderr)
	}
	return exists, nil
}
