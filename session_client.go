package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type SessionStatus struct {
	State int    `json:"state"`
	Msg   string `json:"msg"`
}

type SessionInfo struct {
	Status SessionStatus     `json:"status"`
	Data   map[string]string `json:"data"`
}

func SetSessionByToken(token string) string {
	environment := os.Getenv("GOENV")
	url := "https://xmall.codoon.com/xmall/tokensession?token=%s"

	next_url := fmt.Sprintf(url, token)
	res, err := http.Get(next_url)
	if err != nil {
		fmt.Println("token-session connection error.")
		return ""
	}
	cookie := res.Cookies()
	for _, item := range cookie {
		if item.Name == "sessionid" {
			return item.Value
		}
	}
	return ""
}

func GetUserIdBySession(sessionId string) string {
	environment := os.Getenv("GOENV")
	url := "https://xmall.codoon.com/xmall/get_userid_by_sessionid?sessionid=%s"
	next_url := fmt.Sprintf(url, sessionId)
	res, err := http.Get(next_url)
	if err != nil {
		fmt.Println("token-session connection error.")
		return ""
	}
	defer res.Body.Close()
	ret := SessionInfo{}
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("session_user fetch response failed.")
		return ""
	}
	//	fmt.Println(result)
	err = json.Unmarshal(result, &ret)
	fmt.Println(ret)
	return ret.Data["user_id"]
}
