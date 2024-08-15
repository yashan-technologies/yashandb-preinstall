package checkos

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"preinstall/defines/confdef"
	"preinstall/defines/runtimedef"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/utils/sysutil/osinfoutil"
)

func CheckOS() {
	console.Check("OS操作系统")
	info, err := osinfoutil.GetOSInfo()
	if err != nil {
		console.Fail(fmt.Sprintf("获取操作系统信息失败: %s", err))
		return
	}
	console.OK(fmt.Sprintf("操作系统：%s", info.String()))

	switch runtimedef.GetOSPlatform() {
	case "centos":
		checkCentosVersion(info)
	}

	lang := os.Getenv("LANG")
	console.OK(fmt.Sprintf("语言字符集：%s", lang))
	items := strings.Split(lang, ".")
	if len(items) > 1 {
		charset := confdef.Conf().Limit.Software.Charset
		if items[1] != charset {
			console.Warn(fmt.Sprintf("字符集不为%s，请确认字符集是否符合要求", charset))
		}
	}

	checkTime()
}

func parseCentosVersion(version string) []int64 {
	// 从版本号中解析出主版本号和子版本号，如7.6.1810解析为7.6
	// 根据.分割版本号，取前两位
	versions := strings.Split(version, ".")
	if len(versions) < 2 {
		return nil
	}
	var result []int64
	mainVersion, err := strconv.ParseInt(strings.TrimSpace(versions[0]), 10, 64)
	if err != nil {
		return nil
	}
	result = append(result, mainVersion)
	subVersion, err := strconv.ParseInt(strings.TrimSpace(versions[1]), 10, 64)
	if err != nil {
		return nil
	}
	result = append(result, subVersion)
	return result
}

func checkCentosVersion(info osinfoutil.OSInfo) {
	versions := parseCentosVersion(info.PlatformVersion)
	if versions == nil {
		console.Fail(fmt.Sprintf("解析CentOS版本：%s失败", info.PlatformVersion))
		return
	}
	centosMinVersion := confdef.Conf().Limit.Software.CentosMinVersion
	minVersions := parseCentosVersion(centosMinVersion)
	if minVersions == nil {
		console.Fail(fmt.Sprintf("解析CentOS最低版本：%s失败", centosMinVersion))
		return
	}
	// 对比版本号
	if versions[0] < minVersions[0] ||
		(versions[0] == minVersions[0] && versions[1] < minVersions[1]) {
		console.Fail(fmt.Sprintf("当前CentOS版本：%s，最低要求CentOS版本：%s", info.PlatformVersion, centosMinVersion))
		return
	}
}
