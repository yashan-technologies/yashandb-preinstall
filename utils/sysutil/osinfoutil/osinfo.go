package osinfoutil

import (
	"os"
	"strings"

	"github.com/shirou/gopsutil/v3/host"
)

const (
	UBUNTU_ID = "ubuntu"
	CENTOS_ID = "centos"
	KYLIN_ID  = "kylin"
)

type OSInfo struct {
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platformVersion"`
	KernelArch      string `json:"kernelArch"`
}

func GetOSInfo() (info OSInfo, err error) {
	hostInfo, err := host.Info()
	if err != nil {
		return
	}
	info = OSInfo{
		Platform:        toUpperFirst(hostInfo.Platform),
		PlatformVersion: hostInfo.PlatformVersion,
		KernelArch:      hostInfo.KernelArch,
	}
	if len(info.PlatformVersion) == 0 {
		version, _ := getVersionFromOSRelease()
		info.PlatformVersion = version
	}
	return
}

func (info OSInfo) String() (s string) {
	fields := []string{
		info.Platform,
		info.PlatformVersion,
	}
	return strings.Join(fields, "-") + "." + info.KernelArch
}

func toUpperFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func getVersionFromOSRelease() (version string, err error) {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "VERSION_ID=") {
			version = strings.Split(line, "=")[1]
			version = strings.Trim(version, "\"")
			break
		}
	}
	return
}
