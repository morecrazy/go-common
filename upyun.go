// by liudan,modify by daiping @ 2016-06-28
package common

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"third/upyun"
)

type UpyunParams struct {
	UpBucket   string
	UpUsername string
	UpPassword string
}

var _upyun *upyun.UpYun

func (this *UpyunParams) InitUpyun() {
	if _upyun == nil {
		_upyun = upyun.NewUpYun(this.UpBucket, this.UpUsername, this.UpPassword)
	}

	return
}

func (this *UpyunParams) UploadFile(domain_prefix, key string, value io.Reader, headers map[string]string) (string, http.Header, error) {
	if _upyun == nil {
		return "", nil, errors.New("_upyun not init")
	}

	header, err := _upyun.Put(key, value, true, headers)
	if err != nil {
		return "", nil, err
	} else {
		return this.FileUrl(domain_prefix, key), header, nil
	}

}

func (this *UpyunParams) UploadFileFromUrl(domain_prefix, key, addr string, headers map[string]string) (string, http.Header, error) {
	if _upyun == nil {
		return "", nil, errors.New("_upyun not init")
	}

	_, data, err := SendRequest("GET", addr, nil, nil, nil)
	if err != nil {
		return "", nil, err
	}
	r := strings.NewReader(data)
	return this.UploadFile(domain_prefix, key, r, headers)
}

func (this *UpyunParams) FileUrl(domain_prefix, key string) string {
	//like http://img3.codoon.com/aaaaaaaaa
	return "http://" + domain_prefix + ".codoon.com/" + key
}
