package sysutil

import "os/exec"

func IsCommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
