package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"third/gin"
	"time"

	"testing"
)

const (
	GIN_SERVER_ADDR = "127.0.0.1:18082"
)

func ginServer() {
	engine := gin.New()
	engine.Use(GinKafkaLogger("common-test", "ct", []string{"192.168.1.204:9092"}))
	engine.Use(ReqData2Form())
	engine.POST("/hi", hiHandler)
	go engine.Run(GIN_SERVER_ADDR)
	time.Sleep(2 * time.Second)
}

type HiReq struct {
	UserId string `form:"user_id" binding:"required"`
	Foo    string `form:"foo" binding:"required"`
}

func hiHandler(c *gin.Context) {
	req := &HiReq{}
	if !c.Bind(req) {
		return
	}
	log.Printf("req:%#v", req)
	c.Writer.Write([]byte(req.UserId + req.Foo))
}

func init() {
	ginServer()
}

func TestReqData2Form(t *testing.T) {
	addr := fmt.Sprintf("http://%s/hi", GIN_SERVER_ADDR)
	params := map[string]string{
		"foo": "foo",
	}
	b, _ := json.Marshal(&params)
	// data, err := HttpRequest("POST", addr, params)
	// if string(data) != "foo" {
	// 	t.Fatal("rsp:", string(data), err)
	// }

	request, _ := http.NewRequest("POST", addr, bytes.NewReader(b))
	request.Header.Set("Content-Type", "application/x-www-form-urlencode")
	request.Header.Set("user_id", "abc")
	rsp, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()
	b, _ = ioutil.ReadAll(rsp.Body)
	if string(b) != "abcfoo" {
		t.Fatal("rsp:", string(b), err)
	}
	log.Printf("ret:%s", string(b))

	// time.Sleep(2 * time.Second)
}
