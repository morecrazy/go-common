// by liudan
package common

func Cuts(s string, n int) string {
	if len(s) > n {
		return s[:n]
	} else {
		return s
	}
}

func StrIn(s string, arr []string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}
