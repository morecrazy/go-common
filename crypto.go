// by liudanking 2016.06

package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// I know it is stupid writing a const key in code.
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
