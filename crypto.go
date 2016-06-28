// by liudanking 2016.06

package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math/rand"
	"net/url"
)

// I know it is stupid hard-coding a const key in code.
// We want to introcude hmac into codoon asap and just fire it up.
const _hmac_key = "bes3a3ZnfHzttfkWAUGfxzXPutuRQgUE"

// _hmac_internale_api_token should only be set during init
var _hmac_internale_api_token string

func init() {
	_hmac_internale_api_token = GenMAC("8NK8wjZfJLXtWNUtETPxptNGxcRPFjQw")
}

func GenMAC(msg string) string {
	macF := hmac.New(sha256.New, []byte(_hmac_key))
	macF.Write([]byte(msg))
	mac := macF.Sum(nil)
	return hex.EncodeToString(mac)
}

func VerifyMAC(msg, signature string) bool {
	mac := GenMAC(msg)
	return mac == signature
}

func GenInternalAPIToken() string {
	return _hmac_internale_api_token
}

func VerifyInternalAPIToken(token string) bool {
	return token == _hmac_internale_api_token
}

func SignUrlValue(vm map[string]string) string {
	values := url.Values{}
	for k, v := range vm {
		values.Set(k, v)
	}
	s := values.Encode()
	log.Printf("url encoded:%s", s)
	return GenMAC(s)
}

func VerifyUrlValueSignature(signature string, vm map[string]string) bool {
	return signature == SignUrlValue(vm)
}

func RandNumStr(n int) string {
	dict := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	dictLen := len(dict)
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = dict[rand.Intn(dictLen)]
	}
	return string(b)
}

func RandString(n int) string {
	dict := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	}
	dictLen := len(dict)
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = dict[rand.Intn(dictLen)]
	}
	return string(b)
}
