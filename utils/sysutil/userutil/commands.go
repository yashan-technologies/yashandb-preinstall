package userutil

import "preinstall/defines/bashdef"

var (
	_useraddCommand  = bashdef.CMD_USERADD
	_usermodCommand  = bashdef.CMD_USERMOD
	_groupaddCommand = bashdef.CMD_GROUPADD
)

func SetUseraddCommand(command string) {
	_useraddCommand = command
}

func SetUsermodCommand(command string) {
	_usermodCommand = command
}

func SetGroupaddCommand(command string) {
	_groupaddCommand = command
}
