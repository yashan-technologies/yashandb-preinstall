package fioutil

import (
	"fmt"
	"strings"

	"preinstall/utils/executil"

	"git.yasdb.com/go/yaslog"
)

type Fio struct {
	Cmd       string
	Path      string
	Size      string
	RWMode    string
	BlockSize string
	NumJobs   string
	RunTime   string
	IODepth   string
	Direct    string
}

func NewFio(cmd, path, size, rwMode, blockSize, numJobs, runTime, ioDepth, direct string) Fio {
	return Fio{
		Cmd:       cmd,
		Path:      path,
		Size:      size,
		RWMode:    rwMode,
		BlockSize: blockSize,
		NumJobs:   numJobs,
		RunTime:   runTime,
		IODepth:   ioDepth,
		Direct:    direct,
	}
}

func (f Fio) Run(log yaslog.YasLog) (string, error) {
	execer := executil.NewExecer(log)
	ret, stdout, stderr := execer.Exec(f.Cmd, "--name="+fmt.Sprintf("%s-test", f.RWMode), "--size="+f.Size, "--rw="+f.RWMode,
		"--bs="+f.BlockSize, "--numjobs="+f.NumJobs, "--runtime="+f.RunTime, "--iodepth="+f.IODepth,
		"--direct="+f.Direct, "--filename="+f.Path, "--group_reporting")
	if ret != 0 {
		return "", executil.GenerateError(stdout, stderr)
	}
	return strings.TrimRight(stdout, "\n"), nil
}
