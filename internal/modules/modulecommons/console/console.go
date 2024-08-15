package console

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
)

var (
	Red     = text.Colors{text.FgRed, text.Italic, text.Bold}
	Bule    = text.Colors{text.FgBlue, text.Italic, text.Bold}
	Green   = text.Colors{text.FgGreen, text.Italic, text.Bold}
	Yellow  = text.Colors{text.FgYellow, text.Italic, text.Bold}
	Magenta = text.Colors{text.FgMagenta, text.Italic, text.Bold}
	White   = text.Colors{text.Italic, text.Bold}
)

func Set(msg string) {
	fmt.Println(Magenta.Sprint("[SET]"), White.Sprint(msg))
}

func Info(msg string) {
	fmt.Println(White.Sprint("[INFO]"), msg)
}

func Check(msg string) {
	fmt.Println(Bule.Sprint("[CHECK]"), White.Sprint(msg))
}

func Warn(msg string) {
	fmt.Println(Yellow.Sprint("[WARN]"), msg)
}

func OK(msg string) {
	fmt.Println(Green.Sprint("✔ "), msg)
}

func Fail(msg string) {
	fmt.Println(Red.Sprint("✖ "), msg)
}

func Done() {
	fmt.Println(White.Sprint("[DONE]"))
}
