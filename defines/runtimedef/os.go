package runtimedef

import (
	"strings"

	"preinstall/utils/sysutil/osinfoutil"
)

var (
	_osPlatform string
)

func GetOSPlatform() string {
	return _osPlatform
}

func initOSPlatform() (err error) {
	info, err := osinfoutil.GetOSInfo()
	if err != nil {
		return
	}
	_osPlatform = strings.ToLower(info.Platform)
	return
}
