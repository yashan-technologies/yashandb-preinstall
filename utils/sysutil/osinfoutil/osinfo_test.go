package osinfoutil_test

import (
	"testing"

	"preinstall/utils/sysutil/osinfoutil"
)

func TestGetOSInfo(t *testing.T) {
	info, err := osinfoutil.GetOSInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info.String())
}
