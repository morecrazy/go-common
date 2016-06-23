// by liudan
package common

func Cuts(s string, n int) string {
	if len(s) > n {
		return s[:n]
	} else {
		return s
	}
}
