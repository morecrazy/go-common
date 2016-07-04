// Author 	:	wuql
// Email	: 	wuql@codoon.com
// Date 	: 	2016-7-4
// 给用户发奖章接口

package common

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	GiveMedalUrl = "http://recordmanage.in.codoon.com:4112/medal/give_user_medal"
)

// 奖章元数据
type Medal struct {
	Data struct {
		AcquiredTime     string `json:"acquired_time"`
		Code             string `json:"code"`
		Des              string `json:"des"`
		DisplayGroupShow string `json:"display_group_show"`
		GrayIcon         string `json:"gray_icon"`
		Icon             string `json:"icon"`
		MedalID          int    `json:"medal_id"`
		MiddleGrayIcon   string `json:"middle_gray_icon"`
		MiddleIcon       string `json:"middle_icon"`
		Name             string `json:"name"`
		ShareURL         string `json:"share_url"`
		SmallIcon        string `json:"small_icon"`
	} `json:"data"`
	Status int `json:"status"`
}

// 起时时间2秒
func GiveUserMedal(user_id, code string) (result *Medal, err error) {
	form := url.Values{}
	form.Set("user_id", user_id)
	form.Set("code", code)
	timeout := time.Duration(2 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	request, err := http.NewRequest("POST", GiveMedalUrl, strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(request)
	//
	defer resp.Body.Close()
	//
	if nil != err {
		// add log
		return nil, err
	} else {
		if resp.StatusCode == 200 || resp.StatusCode == 202 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("GiveUserMedal user_id:%s code:%s ioutil.ReadAll error:%v", user_id, code, err)
				return nil, err
			}
			var result *Medal
			if err := json.Unmarshal(body, result); err != nil {
				log.Printf("GiveUserMedal user_id:%s code:%s  json.Unmarshal error:%v", user_id, code, err)
				return nil, err
			} else {
				return result, nil
			}
		} else {
			log.Printf("GiveUserMedal user_id:%s code:%s  resp.StatusCode:%v", user_id, code, resp.StatusCode)
			return nil, errors.New("resp.StatusCode err")
		}
	}
	return
}
