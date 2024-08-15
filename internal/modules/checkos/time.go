package checkos

import (
	"time"

	"preinstall/defines/timedef"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sysutil/osinfoutil"
)

func checkTime() {
	zone, err := osinfoutil.GetTimezone(log.Sugar)
	if err != nil {
		console.Fail("获取时区失败: " + err.Error())
		return
	}
	console.OK("系统时区: " + zone)
	console.OK("本地时间: " + time.Now().Format(timedef.TIME_FORMAT))
	if !osinfoutil.IsTimeUTC8() {
		console.Warn("时区不是UTC+8，请确认时区是否符合要求")
	}
}
