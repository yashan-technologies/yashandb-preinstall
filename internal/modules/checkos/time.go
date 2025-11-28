package checkos

import (
	"fmt"
	"time"

	"preinstall/defines/timedef"
	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/console"
	"preinstall/log"
	"preinstall/utils/sysutil/osinfoutil"
)

func checkTime() {
	zone, err := osinfoutil.GetTimezone(log.Sugar)
	if err != nil {
		console.Fail(fmt.Sprintf(i18n.T("time.get_fail"), err.Error()))
		return
	}
	console.OK(fmt.Sprintf(i18n.T("time.timezone"), zone))
	console.OK(fmt.Sprintf(i18n.T("time.current"), time.Now().Format(timedef.TIME_FORMAT)))
	if !osinfoutil.IsTimeUTC8() {
		console.Warn(i18n.T("time.timezone_warn"))
	}
}
