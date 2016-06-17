// by liudan
package common

import (
	"io"
	"net/http"
	"strings"
	"third/upyun"
)

const (
	UP_BUCKET   = "codoon-img3"
	UP_USERNAME = "codoon"
	UP_PASSWORD = "codoon5401036"
)

var _upyun *upyun.UpYun

func init() {
	if _upyun == nil {
		_upyun = upyun.NewUpYun(UP_BUCKET, UP_USERNAME, UP_PASSWORD)
	}
}

func UploadFile(key string, value io.Reader, headers map[string]string) (string, http.Header, error) {
	header, err := _upyun.Put(key, value, true, headers)
	if err != nil {
		return "", nil, err
	} else {
		return FileUrl(key), header, nil
	}

}

func UploadFileFromUrl(key, addr string, headers map[string]string) (string, http.Header, error) {
	_, data, err := SendRequest("GET", addr, nil, nil, nil)
	if err != nil {
		return "", nil, err
	}
	r := strings.NewReader(data)
	return UploadFile(key, r, headers)
}

func FileUrl(key string) string {
	return "http://img3.codoon.com" + key
}
