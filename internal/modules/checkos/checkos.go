package checkos

import (
	"fmt"

	"preinstall/internal/modules/modulecommons/console"
)

func Check(enableIOTest bool) {
	fmt.Println(console.White.Sprint("<==============================开始检查系统==============================>\n"))
	checkHardWare()
	checkSoftware(enableIOTest)
}

func checkHardWare() {
	CheckOS()
	fmt.Println()

	CheckCPU()
	fmt.Println()

	CheckMemory()
	fmt.Println()

	CheckNetwork()
	fmt.Println()
}

func checkSoftware(enableIOTest bool) {
	CheckFirewall()
	fmt.Println()

	CheckSSH()
	fmt.Println()

	CheckSysctlConf()
	fmt.Println()

	CheckYashanDBUser()
	fmt.Println()

	CheckUlimit()
	fmt.Println()

	CheckYashanDBInstallPath(enableIOTest)
	fmt.Println()
}
