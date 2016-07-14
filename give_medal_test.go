// Author 	:	wuql
// Email	: 	wuql@codoon.com
// Date 	: 	2016-7-4
// 给用户发奖章接口

package common

import (
	"fmt"
	"testing"
)

func TestGiveUserMedal(t *testing.T) {
	user_id := "b44c62d7-cd0b-4027-9e93-ad53163aca84"
	code := "0"
	a, b := GiveUserMedal(user_id, code)
	fmt.Printf("TestGiveUserMedal, a:%v, b:%v ", a, b)
}
