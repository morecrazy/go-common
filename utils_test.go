package common

import "testing"

func TestRandIntRange(t *testing.T) {
	for i := 0; i < 100; i++ {
		a := RandIntRange(1, 2)
		if a > 2 || a < 1 {
			t.Fatalf("TestRandIntRange err a should in range [1,2] not %d", a)
		}
		t.Logf("TestRandIntRange case 1 rand num is:%d", a)
		a = RandIntRange(2, 1)
		if a > 2 || a < 1 {
			t.Fatalf("TestRandIntRange err a should in range [1,2] not %d", a)
		}
		t.Logf("TestRandIntRange case 2 rand num is:%d", a)
	}
}
