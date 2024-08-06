package osinfoutil

import (
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type NetInterface struct {
	Name  string
	Addrs []net.Addr
	Speed int
}

func IsPhysicalNetInterface(interfaceName string) bool {
	devicePath := filepath.Join("/sys/class/net", interfaceName, "device")
	_, err := os.Stat(devicePath)
	return !os.IsNotExist(err)
}

func GetNetInterfaceSpeed(interfaceName string) (int, error) {
	speedFilePath := filepath.Join("/sys/class/net", interfaceName, "speed")
	content, err := os.ReadFile(speedFilePath)
	if err != nil {
		return 0, err
	}
	speedStr := strings.TrimSpace(string(content))
	return strconv.Atoi(speedStr)
}

func GetPhysicalNetInterfaces() ([]NetInterface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var physicalInterfaces []NetInterface
	for _, iface := range interfaces {
		if IsPhysicalNetInterface(iface.Name) {
			addrs, err := iface.Addrs()
			if err != nil {
				return nil, err
			}
			speed, err := GetNetInterfaceSpeed(iface.Name)
			if err != nil {
				return nil, err
			}
			physicalInterfaces = append(physicalInterfaces, NetInterface{
				Name:  iface.Name,
				Addrs: addrs,
				Speed: speed,
			})
		}
	}
	return physicalInterfaces, nil
}
