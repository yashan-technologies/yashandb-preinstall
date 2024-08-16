package std

import (
	"fmt"
	"path"
	"time"

	"preinstall/defines/runtimedef"
	"preinstall/defines/timedef"
	"preinstall/utils/stdutil"
)

const (
	CONSOLE_OUT_FILE = "preinstall.out"
)

var _redirecter *stdutil.Redirecter

func InitRedirecter() error {
	redirecter, err := stdutil.NewRedirecter(genOutput())
	if err != nil {
		return err
	}
	_redirecter = redirecter
	return nil
}

func GetRedirecter() *stdutil.Redirecter {
	return _redirecter
}

func WriteToFile(str string) {
	stdutil.Write(str, _redirecter.GetFileWriter())
}

func WriteToFileAndStdout(str string) {
	stdutil.WriteToStdout(str, _redirecter.GetFileWriter())
}

func genOutput() string {
	return path.Join(runtimedef.ResultPath(), CONSOLE_OUT_FILE+fmt.Sprintf(".%s", time.Now().Format(timedef.TIME_FORMAT_NO_SPACE)))
}
