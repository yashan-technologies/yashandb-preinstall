package checkos

import (
	"fmt"

	"preinstall/i18n"
	"preinstall/internal/modules/modulecommons/console"
)

func Check(enableIOTest bool) {
	fmt.Println(console.White.Sprint(i18n.T("main.start_check") + "\n"))
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
