package setutil_test

import (
	"testing"

	"preinstall/utils/sysutil/setutil"
)

func TestGetTransparentHugePageEnabled(t *testing.T) {
	content, err := setutil.GetTransparentHugePageEnabled()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(content)
}
