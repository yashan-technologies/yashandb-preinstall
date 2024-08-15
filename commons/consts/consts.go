package consts

import (
	"time"

	"preinstall/defines/timedef"
)

const const_backup_suffix = ".preinstall-bak."

func BakupExt() string {
	return const_backup_suffix + time.Now().Format(timedef.TIME_FORMAT_NO_SPACE)
}
