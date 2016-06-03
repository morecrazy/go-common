package common

import (
	"log"
	"time"

	"testing"
)

func TestMAC(t *testing.T) {
	msg := time.Now().String()
	mac := GenMAC(msg)
	if !VerifyMAC(msg, mac) {
		log.Fatal("verify mac failed")
	}

	log.Printf("mac:%s", mac)
}

func TestInternalAPIToken(t *testing.T) {
	token := GenInternalAPIToken()
	if !VerifyInternalAPIToken(token) {
		log.Fatal("verify api internal token failed")
	}

	log.Printf("token:%s", token)
}
