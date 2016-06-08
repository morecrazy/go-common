package common

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestGinStatsdHandler(t *testing.T) {
	addr := fmt.Sprintf("http://%s/hi", GIN_SERVER_ADDR)

	for i := 0; i < 70; i++ {
		log.Printf("request:%d", i)
		data, err := HttpRequest("POST", addr, nil)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Println(string(data))
		}
		time.Sleep(1 * time.Second)
	}

}
