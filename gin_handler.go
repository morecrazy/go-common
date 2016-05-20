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
	"third/kafka"
	"time"
)

// ReqData2Form try to parse request body as json and inject user_id from header to body, if failed, deal with it as form.
// It should be called before your business logic.
func ReqData2Form() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Request.Header.Get("user_id")
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			data, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				log.Printf("read request body error:%v", err)
				return
			}
			v, err := loadJson(bytes.NewReader(data))
			if err != nil {
				// if request data is NOT json format, restore body
				// log.Printf("restore %s to body", string(data))
				values, err := url.ParseQuery(string(data))
				if err != nil {
					log.Printf("parse body data to url values error:%v", err)
					c.Request.Body = ioutil.NopCloser(bytes.NewReader(data))
				} else {
					c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					values.Set("user_id", userId)
					c.Request.Body = ioutil.NopCloser(strings.NewReader(values.Encode()))
				}
			} else {
				c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				v["user_id"] = userId
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

const (
	CODOON_REQUEST_ID = "codoon_request_id"
	CODOON_USER_ID    = "user_id"
	KAFKA_TOPIC       = "codoon-kafka-log"
)

// func GinServiceCoder(code string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if code != "" {
// 			c.Request.Header.Set(CODOON_REQUEST_ID, c.Request.Header.Get(CODOON_REQUEST_ID)+code)
// 		}
// 	}
// }

type KafkaLogger struct {
	RequestId   string
	UserId      string
	ServiceName string
	StartTime   time.Time
	SpendTime   int64
	Method      string
	Host        string
	Api         string
	StatusCode  int
}

func (kl *KafkaLogger) Encode() ([]byte, error) {
	var buf bytes.Buffer
	_, err := fmt.Fprintf(&buf, "%s|%s|%s|%d|%d|%s|%s|%s|%d",
		kl.RequestId,
		kl.UserId,
		kl.ServiceName,
		kl.StartTime.UnixNano()/1e6, // timestamp, ms
		kl.SpendTime,                // ms
		kl.Method,
		kl.Host,
		kl.Api,
		kl.StatusCode,
	)
	if err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func (kl *KafkaLogger) Length() int {
	b, _ := kl.Encode()
	return len(b)
}

// GinKafkaLogger
func GinKafkaLogger(srvName, srvCode string, brockerList []string) gin.HandlerFunc {
	// init producer
	config := kafka.NewConfig()
	config.Producer.RequiredAcks = kafka.WaitForLocal
	config.Producer.Flush.Frequency = 1 * time.Second
	producer, err := kafka.NewAsyncProducer(brockerList, config)
	if err != nil {
		log.Fatalf("create producer error:%v", err)
	}
	inputChannel := producer.Input()
	// monitor kafka error
	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write kafka log entry:", err)
		}
	}()

	return func(c *gin.Context) {
		start := time.Now()

		// colored with current service code
		if srvCode != "" {
			c.Request.Header.Set(CODOON_REQUEST_ID, c.Request.Header.Get(CODOON_REQUEST_ID)+srvCode)
		}

		reqId := c.Request.Header.Get(CODOON_REQUEST_ID)
		userId := c.Request.Header.Get(CODOON_USER_ID)
		method := c.Request.Method
		host := c.Request.Host
		api := c.Request.RequestURI

		c.Next()

		m := &KafkaLogger{
			RequestId:   reqId,
			UserId:      userId,
			ServiceName: srvName,
			StartTime:   start,
			SpendTime:   time.Now().Sub(start).Nanoseconds() / 1e6,
			Method:      method,
			Host:        host,
			Api:         api,
			StatusCode:  c.Writer.Status(),
		}

		inputChannel <- &kafka.ProducerMessage{
			Topic:     KAFKA_TOPIC,
			Partition: 0, // TODO
			Key:       kafka.StringEncoder(srvName),
			Value:     m,
		}

		// log.Printf("kafka msg send:%+v", m)
	}
}
