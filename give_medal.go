// Author 	:	wuql
// Email	: 	wuql@codoon.com
// Date 	: 	2016-7-4
// 给用户发奖章接口

package common

import (
	"encoding/json"
	"errors"
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

func GiveUserMedal(user_id, code string) (*Medal, error) {
	if data, status, err := HttpRequestWithCode(nil, "POST", GiveMedalUrl, map[string]string{"user_id": user_id, "code": code}); err == nil {
		if status == 200 || status == 202 {
			var result Medal
			if err := json.Unmarshal(data, &result); err == nil {
				return &result, nil
			} else {
				return nil, err
			}
		} else {
			return nil, errors.New("status code err")
		}
	} else {
		return nil, err
	}
}
