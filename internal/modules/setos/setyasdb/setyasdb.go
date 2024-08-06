package setyasdb

import (
	"fmt"
	"path"
	"strconv"

	"preinstall/defines/confdef"
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
	console.Set("配置YashanDB安装路径")

	if err := check.CheckRootPrivilege(); err != nil {
		console.Fail("权限不足，配置安装路径，请使用root用户或者sudo执行")
		return
	}

	installPath := confdef.YashanDBConf().InstallPath
	username := confdef.Conf().YashanDBUser.User
	user, err := userutil.GetUserByName(username)
	if err != nil {
		console.Fail(fmt.Sprintf("获取用户信息失败: %s", err))
		return
	}
	uid, err := strconv.Atoi(user.Uid)
	if err != nil {
		console.Fail(fmt.Sprintf("获取用户uid失败: %s", err))
		return
	}
	gid, err := strconv.Atoi(user.Gid)
	if err != nil {
		console.Fail(fmt.Sprintf("获取用户gid失败: %s", err))
		return
	}

	if len(installPath) == 0 {
		console.Warn(fmt.Sprintf("未指定安装路径，使用%s用户的Home目录", username))
		if !userutil.IsUserExists(username) {
			console.Fail(fmt.Sprintf("用户%s不存在，获取Home目录", username))
			return
		}
		installPath = user.HomeDir
	}

	if !fs.IsDirExist(confdef.YashanDBConf().InstallPath) {
		makePath(confdef.YashanDBConf().InstallPath, uid, gid)
		console.OK("创建安装路径：" + confdef.YashanDBConf().InstallPath)
	} else {
		console.OK("安装路径已存在：" + confdef.YashanDBConf().InstallPath)
		if err := fileutil.ReverseChmod(confdef.YashanDBConf().InstallPath, 0755); err != nil {
			console.Fail("修改安装路径权限失败：" + err.Error())
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
			console.Fail(fmt.Sprintf("创建目录%s失败: %s", dir, err))
			return
		}
		console.OK("创建目录：" + dir)
	}

	if setDiskQueneScheduler {
		setScheduler(installPath)
	}
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
		console.Fail(fmt.Sprintf("获取安装路径信息失败: %s", err))
		return
	}
	fname, err := osinfoutil.GetDiskQueneSchedulerPath(info.Device)
	if err != nil {
		console.Fail(fmt.Sprintf("获取磁盘调度器文件失败: %s", err))
		return
	}
	scheduler := confdef.Conf().HostSetting.DiskScheduler
	if err := setutil.SetDiskQueneScheduler(log.Sugar, fname, scheduler); err != nil {
		console.Warn(fmt.Sprintf("设置磁盘调度器失败: %s", err))
	}
	console.OK(fmt.Sprintf("设置%s所在磁盘调度器为：%s", installPath, scheduler))
}
