package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"third/gin"
	"third/httprouter"
	"time"
)

func CodoonGetHeader(c *gin.Context) {
	// 获取token
	r := c.Request
	fmt.Println("+++++++++++header:  ", r.Header)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//io.WriteString(w, "404")
	http.Error(w, "404 page not found", http.StatusNotFound)

}

var sliceOfInts = reflect.TypeOf([]int(nil))
var sliceOfStrings = reflect.TypeOf([]string(nil))

func FormatUserAgent(user_agent string) map[string]interface{} {
	dealed_user_agent := strings.TrimSpace(strings.ToLower(user_agent))
	result := map[string]interface{}{
		"version":          "0.0.0",
		"iner_version":     0,
		"platform":         0,
		"platfrom_version": "",
		"device_type":      "",
	}

	if !strings.Contains(dealed_user_agent, "codoonsport(") {
		return result
	}

	var platform = 1
	if strings.Contains(dealed_user_agent, "ios") {
		platform = 0
	}

	dealed_user_agent = strings.Replace(dealed_user_agent, "codoonsport(", "", -1)
	array_user_agent := strings.Split(dealed_user_agent, ")")
	if array_user_agent == nil {
		return result
	}

	ver_list := strings.Split(array_user_agent[0], ";")
	if len(ver_list) != 3 {
		// 长度有误则为非法agent
		return result
	}
	app_version := strings.Split(ver_list[0], " ")
	platfrom_version := strings.Split(ver_list[1], " ")
	result["version"] = app_version[0]
	result["iner_version"] = app_version[1]
	result["platform"] = platform
	result["platform_version"] = platfrom_version[1]
	result["device_type"] = ver_list[2]
	return result
}

// Compare codoon app version
func CompareVersion(version_a string, version_b string, oper int) (bool, error) {
	// oper:(0, u'>'), (1, u'>='), (2, u'=='), (3, u'<'), (4, u'<='), (other, u'不限制')
	if oper == -1 {
		return true, nil
	}
	int_list_a := []int{}
	int_list_b := []int{}

	if version_a == "" {
		version_a = "0.0.0"
	}

	if version_b == "" {
		version_b = "0.0.0"
	}

	err := StringToIntList(version_a, &int_list_a)
	if err != nil {
		fmt.Errorf("Version format error[version:%v]", version_a)
		return false, err
	}
	for len(int_list_a) < 3 {
		int_list_a = append(int_list_a, 0)
	}
	err = StringToIntList(version_b, &int_list_b)
	if err != nil {
		fmt.Errorf("Version format error[version:%v]", version_b)
		return false, err
	}
	for len(int_list_b) < 3 {
		int_list_b = append(int_list_b, 0)
	}
	if oper == 0 {
		for i := 0; i < len(int_list_a); i++ {
			if int_list_a[i] > int_list_b[i] {
				return true, nil
			} else if int_list_a[i] < int_list_b[i] {
				return false, nil
			} else {
				continue
			}
		}
		return false, nil
	} else if oper == 1 {
		for i := 0; i < 3; i++ {
			if int_list_a[i] > int_list_b[i] {
				return true, nil
			} else if int_list_a[i] < int_list_b[i] {
				return false, nil
			} else {
				continue
			}
		}
		return true, nil
	} else if oper == 2 {
		for i := 0; i < 3; i++ {
			if int_list_a[i] > int_list_b[i] {
				return false, nil
			} else if int_list_a[i] < int_list_b[i] {
				return false, nil
			} else {
				continue
			}
		}
		return true, nil
	} else if oper == 3 {
		for i := 0; i < 3; i++ {
			if int_list_a[i] > int_list_b[i] {
				return false, nil
			} else if int_list_a[i] < int_list_b[i] {
				return true, nil
			} else {
				continue
			}
		}
		return false, nil
	} else if oper == 4 {
		for i := 0; i < 3; i++ {
			if int_list_a[i] > int_list_b[i] {
				return false, nil
			} else if int_list_a[i] < int_list_b[i] {
				return true, nil
			} else {
				continue
			}
		}
		return true, nil
	} else {
		return true, nil
	}
}

func StringToIntList(s string, int_list *[]int) error {
	string_list := strings.Split(s, ".")
	for _, value := range string_list {
		i, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Convert error[%d]", value)
			return err
		}
		*int_list = append(*int_list, i)
	}
	return nil
}

// parse form values to struct via tag.
func ParseForm(form url.Values, obj interface{}) error {
	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	if !IsStructPtr(objT) {
		return fmt.Errorf("%v must be  a struct pointer", obj)
	}
	objT = objT.Elem()
	objV = objV.Elem()

	for i := 0; i < objT.NumField(); i++ {
		fieldV := objV.Field(i)
		if !fieldV.CanSet() {
			continue
		}

		fieldT := objT.Field(i)
		tags := strings.Split(fieldT.Tag.Get("form"), ",")
		var tag string
		if len(tags) == 0 || len(tags[0]) == 0 {
			tag = fieldT.Name
		} else if tags[0] == "-" {
			continue
		} else {
			tag = tags[0]
		}

		value := form.Get(tag)
		if len(value) == 0 {
			continue
		}

		switch fieldT.Type.Kind() {
		case reflect.Bool:
			if strings.ToLower(value) == "on" || strings.ToLower(value) == "1" || strings.ToLower(value) == "yes" {
				fieldV.SetBool(true)
				continue
			}
			if strings.ToLower(value) == "off" || strings.ToLower(value) == "0" || strings.ToLower(value) == "no" {
				fieldV.SetBool(false)
				continue
			}
			b, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			fieldV.SetBool(b)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			x, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			fieldV.SetInt(x)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			x, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			}
			fieldV.SetUint(x)
		case reflect.Float32, reflect.Float64:
			x, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			fieldV.SetFloat(x)
		case reflect.Interface:
			fieldV.Set(reflect.ValueOf(value))
		case reflect.String:
			fieldV.SetString(value)
		case reflect.Struct:
			switch fieldT.Type.String() {
			case "time.Time":
				format := time.RFC3339
				if len(tags) > 1 {
					format = tags[1]
				}
				t, err := time.Parse(format, value)
				if err != nil {
					return err
				}
				fieldV.Set(reflect.ValueOf(t))
			}
		case reflect.Slice:
			if fieldT.Type == sliceOfInts {
				formVals := form[tag]
				fieldV.Set(reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(int(1))), len(formVals), len(formVals)))
				for i := 0; i < len(formVals); i++ {
					val, err := strconv.Atoi(formVals[i])
					if err != nil {
						return err
					}
					fieldV.Index(i).SetInt(int64(val))
				}
			} else if fieldT.Type == sliceOfStrings {
				formVals := form[tag]
				fieldV.Set(reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf("")), len(formVals), len(formVals)))
				for i := 0; i < len(formVals); i++ {
					fieldV.Index(i).SetString(formVals[i])
				}
			}
		}
	}
	return nil
}

func ParseHttpReqToArgs(r *http.Request, args interface{}) error {
	var err error
	ct := r.Header.Get("Content-Type")
	if ct == "application/json" {
		var body []byte
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			Logger.Error("UpdateUserInfo read body err : %s,%v", r.FormValue("user_id"), err)
			return err
		}
		Logger.Debug("body %s", string(body))
		defer r.Body.Close()
		if err := json.Unmarshal(body, args); err != nil {
			Logger.Error("Unmarshal body : %s,%s,%v", r.FormValue("user_id"), string(body), err)
			return err
		}
	} else {
		err = r.ParseForm()
		if nil != err {
			Logger.Error("r.ParseForm err : %v", err)
			err = NewInternalError(DecodeErrCode, err)
		}
		err = ParseForm(r.Form, args)
		if nil != err {
			Logger.Error("ParseForm err : %v", err)
			err = NewInternalError(DecodeErrCode, err)
		}

	}

	return err
}

func WriteRespToBody(w http.ResponseWriter, resp interface{}) error {

	b, err := json.Marshal(resp)
	if err != nil {
		Logger.Error("Marshal json to bytes error :%v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return err
	}
	w.Write(b)
	return err
}

func SendResponse(c *gin.Context, http_code int, data interface{}, err error) error {

	if err != nil {
		CheckError(err)
		c.String(http_code, http.StatusText(http_code))
		return nil
	}

	b, err := json.Marshal(&data)
	if err != nil {
		Logger.Error("Marshal json to bytes error :%v", err)
	}

	//	if len(b) > 10000 {
	//		Logger.Info(string(b[:10000]))
	//	} else {
	//		Logger.Info(string(b))
	//	}

	c.Writer.Write(b)

	return err
}

func SendRequest(http_method, urls string, req_body interface{}, req_form map[string]string, req_raw interface{}) (int, string, error) {
	tr := &http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{Transport: tr}
	form := url.Values{}
	var err error = nil
	var request *http.Request
	var body []byte

	if nil != req_body {
		request, _ = http.NewRequest(http_method, urls, nil)
		request.Header.Set("Content-Type", "application/json")
		b, _ := json.Marshal(req_body)
		request.Body = ioutil.NopCloser(strings.NewReader(string(b)))
	} else if nil != req_form {
		for key, value := range req_form {
			form.Set(key, value)
		}
		if "GET" == http_method {
			request, _ = http.NewRequest(http_method, urls+"?"+form.Encode(), nil)
		} else {
			request, _ = http.NewRequest(http_method, urls, strings.NewReader(form.Encode()))
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else if nil != req_raw {
		var query = []byte(req_raw.(string))
		request, _ = http.NewRequest(http_method, urls, bytes.NewBuffer(query))
		request.Header.Set("Content-Type", "text/plain")
	}

	if nil == req_body && nil == req_form && nil == req_raw {
		request, _ = http.NewRequest(http_method, urls, nil)
	}

	Logger.Debug("request %v", request)

	response, err := client.Do(request)
	if nil != err {
		err = NewInternalError(HttpErrCode, err)
		Logger.Error("send request err :%v", err)
		return http.StatusNotFound, "", err
	}
	// avoid goroutine leak without closing body
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		body, err = ioutil.ReadAll(response.Body)
		if nil == err {
			Logger.Debug("body:%v", string(body))
		}
		return response.StatusCode, string(body), err
	} else {
		err = NewInternalError(HttpErrCode, fmt.Errorf("http code :%d", response.StatusCode))
		Logger.Error("send request err :%v", err)
		return response.StatusCode, "", err
	}

}

func SendFormRequest(http_method, urls string, req_form map[string]string) (int, string, error) {
	return SendRequest(http_method, urls, nil, req_form, nil)
}

func SendJsonRequest(http_method, urls string, req_body interface{}) (int, string, error) {
	return SendRequest(http_method, urls, req_body, nil, nil)
}

func SendRawRequest(http_method, urls string, req_raw interface{}) (int, string, error) {
	return SendRequest(http_method, urls, nil, nil, req_raw)
}

func SMSSendRequest(http_method, urls string, req_body map[string]string) (int, string, error) {

	var err error = nil

	code, result, err := SendFormRequest(http_method, urls, req_body)

	if code == 200 && nil == err {
		result = GetBetweenStr(result, "xmlns=\"http://tempuri.org/\">", "</string>")
	}

	return code, result, err
}

func SendRequestSecure(http_method, urls string, req_form map[string]string, secret string) (int, string, error) {
	tr := &http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{Transport: tr}
	form := url.Values{}
	var err error = nil
	var request *http.Request
	var body []byte

	if nil != req_form {
		for key, value := range req_form {
			form.Set(key, value)
		}
		if "GET" == http_method {
			request, _ = http.NewRequest(http_method, urls+"?"+form.Encode(), nil)
		} else {
			request, _ = http.NewRequest(http_method, urls, strings.NewReader(form.Encode()))
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Set("Authorization", secret)
	}

	Logger.Debug("request %v", request)

	response, err := client.Do(request)
	if nil != err {
		err = NewInternalError(HttpErrCode, err)
		Logger.Error("send request err :%v", err)
		return http.StatusNotFound, "", err
	}

	if response.StatusCode == http.StatusOK {
		defer response.Body.Close()
		body, err = ioutil.ReadAll(response.Body)
		if nil == err {
			Logger.Debug("body:%v", string(body))
		}
		return response.StatusCode, string(body), err
	} else {
		err = NewInternalError(HttpErrCode, fmt.Errorf("http code :%d", response.StatusCode))
		Logger.Error("send request err :%v", err)
		return response.StatusCode, "", err
	}

}

var string_key map[string]int = map[string]int{
	"key":  1,
	"nick": 1,
}

func ForwardHttpToRpc(c *gin.Context, client *RpcClient, method string, args map[string]interface{}, reply interface{}, http_code *int) error {

	r := c.Request
	if nil != c.Request.Body {
		c.ParseBody(&args)
	}

	for _, param := range c.Params {

		value_int, err := strconv.ParseInt(param.Value, 10, 0)
		if err != nil || 1 == string_key[param.Key] {
			args[param.Key] = param.Value
		} else {
			args[param.Key] = value_int
		}
	}

	for key, value := range r.Form {
		value_int, err := strconv.ParseInt(value[0], 10, 0)
		if err != nil || 1 == string_key[key] {
			args[key] = value[0]
		} else {
			args[key] = value_int
		}
	}

	err := client.Call(method, &args, reply)
	if nil != err {
		Logger.Error("Call json rpc error : %s,%v", r.FormValue("user_id"), err)
		*http_code = http.StatusInternalServerError
		return err
	}

	return nil
}

func HttpRequest(method, addr string, params map[string]string) ([]byte, error) {
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	var request *http.Request
	var err error = nil
	if method == "GET" || method == "DELETE" {
		request, err = http.NewRequest(method, addr+"?"+form.Encode(), nil)
		if err != nil {
			return nil, err
		}
	} else {
		request, err = http.NewRequest(method, addr, strings.NewReader(form.Encode()))
		if err != nil {
			return nil, err
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if nil != err {
		log.Printf("httpRequest: Do request (%+v) error:%v", request, err)
		return nil, err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("httpRequest: read response error:%v", err)
		return nil, err
	}
	return data, nil
}
