package runtimedef

import "git.yasdb.com/go/yaserr"

func InitRuntime() error {
	if err := initHomePath(); err != nil {
		return yaserr.Wrapf(err, "init home")
	}
	if err := initExecuteable(); err != nil {
		return yaserr.Wrapf(err, "init executeable")
	}
	if err := initOSPlatform(); err != nil {
		return yaserr.Wrapf(err, "init platform")
	}
	return nil
}
