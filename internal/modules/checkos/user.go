package checkos

import (
	"fmt"
	"strings"

	"preinstall/defines/confdef"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sysutil/userutil"
)

func CheckYashanDBUser() {
	console.Check("数据库用户检查")
	userExists := true
	if !userutil.IsUserExists(confdef.Conf().YashanDBUser.User) {
		userExists = false
		console.Warn(fmt.Sprintf("数据库用户 %s 不存在", confdef.Conf().YashanDBUser.User))
	}

	if !userutil.IsGroupExists(confdef.Conf().YashanDBUser.Group) {
		console.Warn(fmt.Sprintf("数据库用户组 %s 不存在", confdef.Conf().YashanDBUser.Group))
	}
	for _, group := range confdef.Conf().YashanDBUser.AdditionalGroups {
		if !userutil.IsGroupExists(group) {
			console.Warn(fmt.Sprintf("数据库用户附加组 %s 不存在", group))
		}
	}

	if !userExists {
		return
	}
	user, err := userutil.GetUserByName(confdef.Conf().YashanDBUser.User)
	if err != nil {
		console.Fail(fmt.Sprintf("获取用户 %s 信息失败: %s", confdef.Conf().YashanDBUser.User, err.Error()))
		return
	}
	primaryGroup, err := userutil.GetGroupByID(user.Gid)
	if err != nil {
		console.OK("数据库用户：" + user.Username)
		console.OK("家目录：" + user.HomeDir)
		console.Fail(fmt.Sprintf("获取用户 %s 主组信息失败: %s", confdef.Conf().YashanDBUser.User, err.Error()))
		return
	}

	defer func() {
		if primaryGroup.Name != confdef.Conf().YashanDBUser.Group {
			console.Warn(fmt.Sprintf("数据库用户 %s 主用户组 %s 与配置文件中的主用户组 %s 不一致", user.Username, primaryGroup.Name, confdef.Conf().YashanDBUser.Group))
		}
		if !userutil.CheckSudoForUser(log.Sugar, user.Username) {
			console.Warn(fmt.Sprintf("数据库用户 %s 无sudo权限，需要手动配置 /etc/sudoers", user.Username))
		}
	}()

	gids, err := user.GroupIds()
	if err != nil {
		console.OK("数据库用户：" + user.Username)
		console.OK("主用户组：" + primaryGroup.Name)
		console.OK("家目录：" + user.HomeDir)
		console.Fail(fmt.Sprintf("获取用户 %s 附加用户组信息失败: %s", user.Username, err.Error()))
		return
	}

	var names []string
	nameMap := make(map[string]struct{})
	for _, gid := range gids {
		if gid == user.Gid {
			continue
		}
		group, err := userutil.GetGroupByID(gid)
		if err != nil {
			console.OK("数据库用户：" + user.Username)
			console.OK("主用户组：" + primaryGroup.Name)
			console.OK("家目录：" + user.HomeDir)
			console.Fail(fmt.Sprintf("获取用户 %s 附加组信息失败: %s", user.Username, err.Error()))
			return
		}
		names = append(names, group.Name)
		nameMap[group.Name] = struct{}{}
	}
	console.OK("数据库用户：" + user.Username)
	console.OK("主用户组：" + primaryGroup.Name)
	console.OK("附加用户组：" + strings.Join(names, ", "))
	console.OK("家目录：" + user.HomeDir)

	for _, group := range confdef.Conf().YashanDBUser.AdditionalGroups {
		if _, ok := nameMap[group]; !ok {
			console.Warn(fmt.Sprintf("数据库用户 %s 未被加入附加组 %s", user.Username, group))
		}
	}
}
