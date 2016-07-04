// by liudan

package common

import (
	"regexp"
)

var _reg_digit = regexp.MustCompile(`\d+`)

func IsDigital(s string) bool {
	return _reg_digit.MatchString(s)
}
