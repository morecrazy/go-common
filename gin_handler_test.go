package common

import (
	"fmt"
	"log"
	"third/gin"

	"testing"
)

const (
	GIN_SERVER_ADDR = "127.0.0.1:18082"
)

func ginServer() {
	engine := gin.New()
	engine.Use(ReqData2Form())
	engine.POST("/hi", hiHandler)
	go engine.Run(GIN_SERVER_ADDR)
}

type HiReq struct {
	Foo string `form:"foo"`
}

func hiHandler(c *gin.Context) {
	req := &HiReq{}
	if !c.Bind(req) {
		return
	}
	log.Printf("req:%#v", req)
	c.Writer.Write([]byte(req.Foo))
}

func init() {
	ginServer()
}

func TestReqData2Form(t *testing.T) {
	addr := fmt.Sprintf("http://%s/hi", GIN_SERVER_ADDR)
	params := map[string]string{
		"foo": "foo",
	}
	data, err := HttpRequest("POST", addr, params)
	if string(data) != "foo" {
		t.Fatal("rsp:", string(data), err)
	}
}
