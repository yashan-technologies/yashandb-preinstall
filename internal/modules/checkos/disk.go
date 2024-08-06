package checkos

import (
	"fmt"
	"path"
	"strings"

	"preinstall/defines/confdef"
	"preinstall/defines/runtimedef"
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
	console.Check("数据库安装路径")
	installPath := confdef.YashanDBConf().InstallPath
	username := confdef.Conf().YashanDBUser.User
	if len(installPath) == 0 {
		console.Warn(fmt.Sprintf("未指定安装路径，使用%s用户的Home目录", username))
		if !userutil.IsUserExists(username) {
			console.Fail(fmt.Sprintf("用户%s不存在，无法获取Home目录", username))
			return
		}
		user, err := userutil.GetUserByName(username)
		if err != nil {
			console.Fail(fmt.Sprintf("获取用户信息失败: %s", err))
			return
		}
		installPath = user.HomeDir
	}

	info, err := osinfoutil.GetDiskInfoByPath(installPath)
	if err != nil {
		console.Fail(fmt.Sprintf("获取安装路径信息失败: %s", err))
		return
	}

	var checkFreeSpace bool
	usage, err := disk.Usage(info.Mountpoint)
	if err != nil {
		console.OK("安装路径：" + installPath)
		console.OK("磁盘/分区：" + info.Device)
		console.OK("文件系统：" + info.Fstype)
		console.Fail(fmt.Sprintf("获取安装路径的可用空间失败: %s", err))
	} else {
		checkFreeSpace = true
		console.OK("安装路径：" + installPath)
		console.OK("磁盘/分区：" + info.Device)
		console.OK("文件系统：" + info.Fstype)
		console.OK("可用空间：" + size.GenHumanReadableSize(float64(usage.Free), 2))
	}

	var fsTypeMatched bool
	for _, t := range confdef.Conf().Limit.Hardware.InstallPathFsTypes {
		if t == info.Fstype {
			fsTypeMatched = true
			break
		}
	}
	if !fsTypeMatched {
		console.Warn(fmt.Sprintf("安装路径的文件系统为：%s，建议您使用 %s 类型的文件系统", info.Fstype,
			console.Green.Sprint(strings.Join(confdef.Conf().Limit.Hardware.InstallPathFsTypes, "、"))))
	}

	if checkFreeSpace {
		minFreeSpace := uint64(confdef.Conf().Limit.Hardware.InstallPathMinFreeGB) * 1024 * 1024 * 1024
		if usage.Free < minFreeSpace {
			console.Warn(fmt.Sprintf("安装路径的磁盘可用空间为：%s，建议您保留 %s 及以上的可用空间",
				size.GenHumanReadableSize(float64(usage.Free), 2),
				console.Green.Sprintf("%dG", confdef.Conf().Limit.Hardware.InstallPathMinFreeGB)))
			check.AddCheckCount()
		}
	}

	if !fs.IsDirExist(installPath) {
		console.Warn(fmt.Sprintf("安装路径%s不存在或无权限访问，请检查", installPath))
	}

	fname, err := osinfoutil.GetDiskQueneSchedulerPath(info.Device)
	if err != nil {
		console.Fail(fmt.Sprintf("获取磁盘调度器文件失败: %s", err))
		return
	}

	content, err := osinfoutil.GetDiskQueneScheduler(fname)
	if err != nil {
		console.Fail(fmt.Sprintf("获取磁盘调度器失败: %s", err))
		return
	}

	console.OK(fmt.Sprintf("磁盘调度器[cat %s]：%s", fname, content))

	runIOTest(enableIOTest, path.Join(installPath, "fio.test"))
}

func runIOTest(enableIOTest bool, path string) {
	if !enableIOTest {
		return
	}
	cmd := runtimedef.Fio()
	c := confdef.Conf().Fio
	for _, mode := range c.RWMode {
		console.Info(fmt.Sprintf("开始进行%s模式的磁盘I/O性能测试......", mode))
		fio := fioutil.NewFio(cmd, path, c.Size, mode, c.BlockSize, c.NumJobs, c.RunTime, c.IODepth, c.Direct)
		out, err := fio.Run(log.Sugar)
		if err != nil {
			console.Fail(fmt.Sprintf("%s模式磁盘I/O性能测试失败: \n%s", mode, err))
			return
		}
		console.OK(fmt.Sprintf("%s模式磁盘I/O性能测试结果: %s", mode, out))
	}
}
