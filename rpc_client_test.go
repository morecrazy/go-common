package common

import (
	. "backend/common/protocol"
	"fmt"
	"testing"
)

func TestGetProfileById(t *testing.T) {
	userId := "83845d28-171c-4a20-9053-69463eea667b"
	user, err := GetProfileById(userId)
	if err != nil {
		fmt.Println("%s,%v", userId, err)
		t.Fatalf("error")
	}
}
