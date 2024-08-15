package errdef

import "errors"

var ErrNeedRootPrivilege = errors.New("root or sudo privilege required")
