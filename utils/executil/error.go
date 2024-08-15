package executil

import (
	"fmt"
	"strings"
)

func GenerateError(stdout, stderr string) error {
	return fmt.Errorf(strings.ReplaceAll(strings.Join([]string{stdout, stderr}, " "), "\n", ""))
}
