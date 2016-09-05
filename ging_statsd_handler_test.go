package common

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestGinStatsdHandler(t *testing.T) {
	addr := fmt.Sprintf("http://%s/hi", GIN_SERVER_ADDR)
	params := map[string]string{
		"user_id": "abc",
		"foo":     "foo",
	}
	for i := 0; i < 70; i++ {
		log.Printf("request:%d", i)
		data, err := HttpRequest("POST", addr, params)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Println(string(data))
		}
		time.Sleep(1 * time.Second)
	}

}
