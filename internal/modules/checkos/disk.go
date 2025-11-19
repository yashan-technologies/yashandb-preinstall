package checkos

import (
	"fmt"
	"path"
	"strings"

	"preinstall/defines/confdef"
	"preinstall/defines/runtimedef"
	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/check"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/fioutil"
	"preinstall/utils/sysutil/osinfoutil"
	"preinstall/utils/sysutil/userutil"

	"git.yasdb.com/go/yasutil/fs"
	"git.yasdb.com/go/yasutil/size"
	"github.com/shirou/gopsutil/v3/disk"
)

func CheckYashanDBInstallPath(enableIOTest bool) {
	console.Check(i18n.T("disk.check"))
	installPath := confdef.YashanDBConf().InstallPath
	username := confdef.Conf().YashanDBUser.User
	if len(installPath) == 0 {
		console.Warn(fmt.Sprintf(i18n.T("set.yasdb_path.use_home"), username))
		if !userutil.IsUserExists(username) {
			console.Fail(fmt.Sprintf(i18n.T("set.yasdb_path.user_not_exist"), username))
			return
		}
		user, err := userutil.GetUserByName(username)
		if err != nil {
			console.Fail(fmt.Sprintf(i18n.T("set.yasdb_path.get_user_fail"), err))
			return
		}
		installPath = user.HomeDir
	}

	info, err := osinfoutil.GetDiskInfoByPath(installPath)
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("disk.get_info_fail"), err))
		return
	}

	var checkFreeSpace bool
	usage, err := disk.Usage(info.Mountpoint)
	if err != nil {
		console.OK(fmt.Sprintf(i18n.T("disk.install_path"), installPath))
		console.OK(fmt.Sprintf(i18n.T("disk.device"), info.Device))
		console.OK(fmt.Sprintf(i18n.T("disk.fs_type"), info.Fstype))
		console.Fail(fmt.Sprintf(i18n.T("disk.get_info_fail"), err))
	} else {
		checkFreeSpace = true
		console.OK(fmt.Sprintf(i18n.T("disk.install_path"), installPath))
		console.OK(fmt.Sprintf(i18n.T("disk.device"), info.Device))
		console.OK(fmt.Sprintf(i18n.T("disk.fs_type"), info.Fstype))
		console.OK(fmt.Sprintf(i18n.T("disk.free"), size.GenHumanReadableSize(float64(usage.Free), 2)))
	}

	var fsTypeMatched bool
	for _, t := range confdef.Conf().Limit.Hardware.InstallPathFsTypes {
		if t == info.Fstype {
			fsTypeMatched = true
			break
		}
	}
	if !fsTypeMatched {
		console.Warn(fmt.Sprintf(i18n.T("disk.fs_type_warn"), info.Fstype,
			strings.Join(confdef.Conf().Limit.Hardware.InstallPathFsTypes, ", ")))
	}

	if checkFreeSpace {
		minFreeSpace := uint64(confdef.Conf().Limit.Hardware.InstallPathMinFreeGB) * 1024 * 1024 * 1024
		if usage.Free < minFreeSpace {
			console.Warn(fmt.Sprintf(i18n.T("disk.free_space_warn"),
				size.GenHumanReadableSize(float64(usage.Free), 2),
				fmt.Sprintf("%dG", confdef.Conf().Limit.Hardware.InstallPathMinFreeGB)))
			check.AddCheckCount()
		}
	}

	if !fs.IsDirExist(installPath) {
		console.Warn(fmt.Sprintf(i18n.T("disk.path_not_exist"), installPath))
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
		if !fs.IsDirExist(dir) {
			console.Warn(fmt.Sprintf(i18n.T("disk.dir_not_exist"), dir))
		}
	}

	fname, err := osinfoutil.GetDiskQueneSchedulerPath(info.Device)
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("set.yasdb_path.get_scheduler_fail"), err))
		return
	}

	content, err := osinfoutil.GetDiskQueneScheduler(fname)
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("disk.get_scheduler_fail"), err))
		return
	}

	console.OK(fmt.Sprintf(i18n.T("disk.scheduler"), content))

	runIOTest(enableIOTest, path.Join(installPath, "fiotest.preinstall"))
}

func runIOTest(enableIOTest bool, filename string) {
	if !enableIOTest {
		return
	}
	console.Info(fmt.Sprintf(i18n.T("disk.io_testing"), filename))
	cmd := runtimedef.Fio()
	c := confdef.Conf().Fio
	for _, mode := range c.RWMode {
		console.Info(fmt.Sprintf(i18n.T("disk.io_mode_testing"), mode))
		fio := fioutil.NewFio(cmd, filename, c.Size, mode, c.BlockSize, c.NumJobs, c.RunTime, c.IODepth, c.Direct)
		out, err := fio.Run(log.Sugar)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "permission denied") {
				console.Fail(fmt.Sprintf(i18n.T("disk.io_permission_fail"), mode, filename))
				return
			}
			console.Fail(fmt.Sprintf(i18n.T("disk.io_fail"), mode, err))
			return
		}
		console.OK(fmt.Sprintf(i18n.T("disk.io_result"), mode, out))
	}
}
