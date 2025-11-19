package setyasdb

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"preinstall/defines/confdef"
	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/check"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sysutil/fileutil"
	"preinstall/utils/sysutil/osinfoutil"
	"preinstall/utils/sysutil/setutil"
	"preinstall/utils/sysutil/userutil"

	"git.yasdb.com/go/yasutil/fs"
)

func SetYashanDBPath(setDiskQueneScheduler bool) {
	console.Set(i18n.T("set.yasdb_path.title"))
	console.Info(fmt.Sprintf(i18n.T("set.yasdb_path.info"), confdef.YashanDBConf().InstallPath))

	if err := check.CheckRootPrivilege(); err != nil {
		console.Fail(i18n.T("set.yasdb_path.no_privilege"))
		return
	}

	installPath := confdef.YashanDBConf().InstallPath
	username := confdef.Conf().YashanDBUser.User
	user, err := userutil.GetUserByName(username)
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("set.yasdb_path.get_user_fail"), err))
		return
	}
	uid, err := strconv.Atoi(user.Uid)
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("set.yasdb_path.get_uid_fail"), err))
		return
	}
	gid, err := strconv.Atoi(user.Gid)
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("set.yasdb_path.get_gid_fail"), err))
		return
	}

	if len(installPath) == 0 {
		console.Warn(fmt.Sprintf(i18n.T("set.yasdb_path.use_home"), username))
		if !userutil.IsUserExists(username) {
			console.Fail(fmt.Sprintf(i18n.T("set.yasdb_path.user_not_exist"), username))
			return
		}
		installPath = user.HomeDir
	}

	if !fs.IsDirExist(confdef.YashanDBConf().InstallPath) {
		makePath(confdef.YashanDBConf().InstallPath, uid, gid)
		console.OK(fmt.Sprintf(i18n.T("set.yasdb_path.created"), confdef.YashanDBConf().InstallPath))
	} else {
		console.OK(fmt.Sprintf(i18n.T("set.yasdb_path.exist"), confdef.YashanDBConf().InstallPath))
		if err := fileutil.ReverseChmod(confdef.YashanDBConf().InstallPath, 0755); err != nil {
			console.Fail(fmt.Sprintf(i18n.T("set.yasdb_path.chmod_fail"), err.Error()))
			return
		}
	}

	yasdbHome := confdef.YashanDBConf().YasdbHome
	yasdbData := confdef.YashanDBConf().YasdbData
	yasdbBack := confdef.YashanDBConf().YasdbBack
	if !path.IsAbs(yasdbHome) {
		yasdbHome = path.Join(installPath, yasdbHome)
	}
	if !path.IsAbs(yasdbData) {
		yasdbData = path.Join(installPath, yasdbData)
	}
	if !path.IsAbs(yasdbBack) {
		yasdbBack = path.Join(installPath, yasdbBack)
	}
	dirs := []string{yasdbHome, yasdbData, yasdbBack}
	for _, dir := range confdef.YashanDBConf().YasdbBackSubdirs {
		dirs = append(dirs, path.Join(yasdbBack, dir))
	}
	for _, dir := range dirs {
		if err := makePath(dir, uid, gid); err != nil {
			console.Fail(fmt.Sprintf(i18n.T("set.yasdb_path.create_dir_fail"), dir, err))
			return
		}
		console.OK(fmt.Sprintf(i18n.T("set.yasdb_path.dir_created"), dir))
	}

	if setDiskQueneScheduler {
		setScheduler(installPath)
	}
	console.Done()
}

func makePath(path string, uid, gid int) error {
	if !fs.IsDirExist(path) {
		if err := fs.Mkdir(path); err != nil {
			return err
		}
	}
	if err := fileutil.ReverseChmod(path, 0755); err != nil {
		return err
	}
	if err := fileutil.RecursiveChown(path, uid, gid); err != nil {
		return err
	}
	return nil
}

func setScheduler(installPath string) {
	info, err := osinfoutil.GetDiskInfoByPath(installPath)
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("set.yasdb_path.get_disk_fail"), err))
		return
	}
	fname, err := osinfoutil.GetDiskQueneSchedulerPath(info.Device)
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("set.yasdb_path.get_scheduler_fail"), err))
		return
	}
	scheduler := confdef.Conf().HostSetting.DiskScheduler
	if err := setutil.SetDiskQueneScheduler(log.Sugar, fname, scheduler); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "invalid argument") {
			console.Fail(fmt.Sprintf(i18n.T("set.yasdb_path.scheduler_not_support"), scheduler))
			return
		}
		console.Fail(fmt.Sprintf(i18n.T("set.yasdb_path.set_scheduler_fail"), err))
		return
	}
	console.OK(fmt.Sprintf(i18n.T("set.yasdb_path.scheduler_set"), installPath, scheduler))
}
