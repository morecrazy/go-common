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

func TestSignUrlValue(t *testing.T) {
	vm := map[string]string{
		"user_id": "abcd",
		"content": "测试中文",
		"a":       "b",
	}
	signature := SignUrlValue(vm)
	log.Printf("signature:%s", signature)
	if !VerifyUrlValueSignature(signature, vm) {
		log.Fatal("verify url value signature failed")
	}
}
