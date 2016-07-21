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

func ginTestHandler(c *gin.Context) {
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencode")
	c.Request.Header.Set("Content-Length", "0")
}

func ginServer() {
	engine := gin.New()
	engine.Use(ginTestHandler)
	engine.Use(GinStatter("120.26.81.104:8125", "common-test"))
	// engine.Use(GinKafkaLogger("common-test", "ct", []string{"192.168.1.204:9092"}))
	engine.Use(ReqData2Form())
	engine.POST("/hi", hiHandler)
	engine.POST("/uid", uidHandler)
	engine.POST("/json", reqJsonHanlder)
	go engine.Run(GIN_SERVER_ADDR)
	time.Sleep(2 * time.Second)
}

type HiReq struct {
	UserId string `form:"user_id" binding:"required"`
	Foo    string `form:"foo" binding:"required"`
}

func hiHandler(c *gin.Context) {
	log.Printf("request header:%+v", c.Request.Header)
	req := &HiReq{}
	if !c.Bind(req) {
		c.JSON(http.StatusOK, gin.H{"ret": "bind failed"})
		return
	}
	log.Printf("req:%#v", req)
	log.Printf("req_user_id:%s", c.DefaultPostForm("req_user_id", "defaultValue"))
	c.Writer.Write([]byte(req.UserId + req.Foo))
}

func reqJsonHanlder(c *gin.Context) {
	expect := map[string]string{
		"string": "test",
		"pics":   `[{"url":"http://img3.codoon.com/f2cdae30-a"}]`,
		"dict":   `{"label_content":"sport"}`,
		"float":  "10.2",
		"int":    "20",
	}

	for k, v1 := range expect {
		v2 := c.PostForm(k)
		if v1 != v2 {
			log.Printf("[%s]!=[%s]", v1, v2)
			c.Writer.Write([]byte("Error"))
			return
		}
	}

	c.Writer.Write([]byte("OK"))

}

func uidHandler(c *gin.Context) {
	c.Writer.Write([]byte(c.DefaultPostForm("user_id", "")))
}

func init() {
	ginServer()
}

func TestReqData2Form(t *testing.T) {
	// test case 1: normal json request
	addr := fmt.Sprintf("http://%s/hi", GIN_SERVER_ADDR)
	params := map[string]string{
		"foo": "foo",
	}
	b, _ := json.Marshal(&params)

	request, _ := http.NewRequest("POST", addr, bytes.NewReader(b))
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

	// test case 2: empty body request
	addr = fmt.Sprintf("http://%s/uid", GIN_SERVER_ADDR)
	request, _ = http.NewRequest("POST", addr, nil)
	request.Header.Set("user_id", "abc")
	rsp, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()
	b, _ = ioutil.ReadAll(rsp.Body)
	if string(b) != "abc" {
		t.Fatal("rsp:", string(b), err)
	}
	log.Printf("ret:%s", string(b))

	// test case 3: json request with user_id as key
	addr = fmt.Sprintf("http://%s/hi", GIN_SERVER_ADDR)
	params = map[string]string{
		"user_id": "test_user_id",
		"foo":     "foo",
	}
	b, _ = json.Marshal(&params)
	request, _ = http.NewRequest("POST", addr, bytes.NewReader(b))
	request.Header.Set("user_id", "abc")
	rsp, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()
	b, _ = ioutil.ReadAll(rsp.Body)
	if string(b) != "abcfoo" {
		t.Fatal("rsp:", string(b), err)
	}
	log.Printf("ret:%s", string(b))
}

func TestJsonReq(t *testing.T) {
	addr := fmt.Sprintf("http://%s/json", GIN_SERVER_ADDR)
	var pics interface{}
	var dict interface{}
	json.Unmarshal([]byte(`[{"url":"http://img3.codoon.com/f2cdae30-a"}]`), &pics)
	json.Unmarshal([]byte(`{"label_content":"sport"}`), &dict)

	params := map[string]interface{}{
		"string": "test",
		"pics":   pics,
		"dict":   dict,
		"float":  10.2,
		"int":    20,
	}

	b, _ := json.Marshal(&params)
	request, _ := http.NewRequest("POST", addr, bytes.NewReader(b))
	rsp, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()
	b, _ = ioutil.ReadAll(rsp.Body)
	if string(b) != "OK" {
		t.Fatal("rsp:", string(b), err)
	}
	log.Printf("ret:%s", string(b))

}
