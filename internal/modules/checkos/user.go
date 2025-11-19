package checkos

import (
	"fmt"
	"strings"

	"preinstall/defines/confdef"
	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sysutil/userutil"
)

func CheckYashanDBUser() {
	console.Check(i18n.T("user.check"))
	userExists := true
	if !userutil.IsUserExists(confdef.Conf().YashanDBUser.User) {
		userExists = false
		console.Warn(fmt.Sprintf(i18n.T("user.not_exist"), confdef.Conf().YashanDBUser.User))
	}

	if !userutil.IsGroupExists(confdef.Conf().YashanDBUser.Group) {
		console.Warn(fmt.Sprintf(i18n.T("user.group_not_exist"), confdef.Conf().YashanDBUser.Group))
	}
	for _, group := range confdef.Conf().YashanDBUser.AdditionalGroups {
		if !userutil.IsGroupExists(group) {
			console.Warn(fmt.Sprintf(i18n.T("user.additional_group_not_exist"), group))
		}
	}

	if !userExists {
		return
	}
	user, err := userutil.GetUserByName(confdef.Conf().YashanDBUser.User)
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("user.get_info_fail"), confdef.Conf().YashanDBUser.User, err.Error()))
		return
	}
	primaryGroup, err := userutil.GetGroupByID(user.Gid)
	if err != nil {
		console.OK(fmt.Sprintf(i18n.T("user.exist"), user.Username))
		console.OK(fmt.Sprintf(i18n.T("user.home"), user.HomeDir))
		console.Fail(fmt.Sprintf(i18n.T("user.get_primary_group_fail"), confdef.Conf().YashanDBUser.User, err.Error()))
		return
	}

	defer func() {
		if primaryGroup.Name != confdef.Conf().YashanDBUser.Group {
			console.Warn(fmt.Sprintf(i18n.T("user.primary_group_mismatch"), user.Username, primaryGroup.Name, confdef.Conf().YashanDBUser.Group))
		}
		if !userutil.CheckSudoForUser(log.Sugar, user.Username) {
			console.Warn(fmt.Sprintf(i18n.T("user.no_sudo"), user.Username))
		} else {
			console.OK(i18n.T("user.sudo_configured"))
		}
	}()

	gids, err := user.GroupIds()
	if err != nil {
		console.OK(fmt.Sprintf(i18n.T("user.exist"), user.Username))
		console.OK(fmt.Sprintf(i18n.T("user.primary_group"), primaryGroup.Name))
		console.OK(fmt.Sprintf(i18n.T("user.home"), user.HomeDir))
		console.Fail(fmt.Sprintf(i18n.T("user.get_additional_groups_fail"), user.Username, err.Error()))
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
			console.OK(fmt.Sprintf(i18n.T("user.exist"), user.Username))
			console.OK(fmt.Sprintf(i18n.T("user.primary_group"), primaryGroup.Name))
			console.OK(fmt.Sprintf(i18n.T("user.home"), user.HomeDir))
			console.Fail(fmt.Sprintf(i18n.T("user.get_group_info_fail"), user.Username, err.Error()))
			return
		}
		names = append(names, group.Name)
		nameMap[group.Name] = struct{}{}
	}
	console.OK(fmt.Sprintf(i18n.T("user.exist"), user.Username))
	console.OK(fmt.Sprintf(i18n.T("user.primary_group"), primaryGroup.Name))
	console.OK(fmt.Sprintf(i18n.T("user.additional_groups"), strings.Join(names, ", ")))
	console.OK(fmt.Sprintf(i18n.T("user.home"), user.HomeDir))

	for _, group := range confdef.Conf().YashanDBUser.AdditionalGroups {
		if _, ok := nameMap[group]; !ok {
			console.Warn(fmt.Sprintf(i18n.T("user.not_in_group"), user.Username, group))
		}
	}
}
