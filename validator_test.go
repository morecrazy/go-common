package common

import "testing"

func TestIsDigit(t *testing.T) {
	s := "1234567890"
	if !IsDigital(s) {
		t.Fatalf("%s is digit", s)
	}
}
