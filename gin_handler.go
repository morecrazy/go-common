package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"third/gin"
	"time"
)

// ReqData2Form try to parse request body as json, if failed, deal with it as form.
// It should be called before your business logic.
func ReqData2Form() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			data, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				log.Printf("read request body error:%v", err)
				return
			}
			v, err := loadJson(bytes.NewReader(data))
			if err != nil {
				// if not request data is NOT json format, restore body
				// log.Printf("restore %s to body", string(data))
				c.Request.Body = ioutil.NopCloser(bytes.NewReader(data))

			} else {
				form := map2Form(v)
				c.Request.Body = ioutil.NopCloser(strings.NewReader(form.Encode()))
			}
		}
	}
}

func loadJson(r io.Reader) (map[string]interface{}, error) {
	decoder := json.NewDecoder(r)
	decoder.UseNumber()
	var v map[string]interface{}
	err := decoder.Decode(&v)
	if err != nil {
		// log.Printf("loadJson decode error:%v", err)
		return nil, err
	}
	return v, nil
}

func map2Form(v map[string]interface{}) url.Values {
	form := url.Values{}
	var vStr string
	var ok bool
	for key, value := range v {
		if vStr, ok = value.(string); !ok {
			vStr = fmt.Sprintf("%v", value)
		}
		form.Set(key, vStr)
	}
	return form
}

type SlowLogger interface {
	Notice(format string, params ...interface{})
	Warning(format string, params ...interface{})
}

func GinSlowLogger(slog SlowLogger, threshold time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		startAt := time.Now()

		c.Next()

		endAt := time.Now()
		latency := endAt.Sub(startAt)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		if latency > threshold {
			slog.Warning("[GIN Slowlog] %v | %3d | %12v | %s | %-7s %s %s\n%s",
				endAt.Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				c.Request.URL.String(),
				c.Request.URL.Opaque,
				c.Errors.String())
		}
	}
}
