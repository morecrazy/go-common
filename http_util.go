package common

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
		err := r.ParseForm()
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

func SendRequest(http_method, urls string, req_body interface{}, req_form map[string]string) (int, string, error) {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
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
	}

	Logger.Debug("request %v", request)

	response, err := client.Do(request)
	if nil != err {
		err = NewInternalError(HttpErrCode, err)
		Logger.Error("send request err :%v", err)
		return 200, "", err
	}

	if response.StatusCode == 200 {
		defer response.Body.Close()
		body, err = ioutil.ReadAll(response.Body)
		if nil == err {
			Logger.Debug("body:%v", string(body))
		}
	} else {
		err = NewInternalError(HttpErrCode, fmt.Errorf("http code :%d", response.StatusCode))
		Logger.Error("send request err :%v", err)
		return 200, "", err
	}

	return response.StatusCode, string(body), err
}

func SendFormRequest(http_method, urls string, req_form map[string]string) (int, string, error) {
	return SendRequest(http_method, urls, nil, req_form)
}

func SendJsonRequest(http_method, urls string, req_body interface{}) (int, string, error) {
	return SendRequest(http_method, urls, req_body, nil)
}

func SMSSendRequest(http_method, urls string, req_body map[string]string) (int, string, error) {

	var err error = nil

	code, result, err := SendFormRequest(http_method, urls, req_body)

	if code == 200 && nil == err {
		result = GetBetweenStr(result, "xmlns=\"http://tempuri.org/\">", "</string>")
	}

	return code, result, err
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

func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			err := recover()
			if err != nil {
				switch err.(type) {
				case error:
					CheckError(err.(error))
				default:
					err := errors.New(fmt.Sprint(err))
					CheckError(err)
				}

				stack := stack(3)
				Logger.Error("PANIC: %s\n%s", err, stack)

				c.Writer.WriteHeader(http.StatusInternalServerError)
			}

		}()

		c.Next()
	}
}

func MyRecovery() {

	err := recover()
	if err != nil {
		switch err.(type) {
		case error:
			CheckError(err.(error))
		default:
			err := errors.New(fmt.Sprint(err))
			CheckError(err)
		}

		stack := stack(3)
		Logger.Error("PANIC: %s\n%s", err, stack)
	}

}

func GinLogger() gin.HandlerFunc {

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		Logger.Notice("[GIN] %v | %3d | %12v | %s | %-7s %s %s\n%s",
			end.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			c.Request.URL.String(),
			c.Request.URL.Opaque,
			c.Errors.String())

	}
}
