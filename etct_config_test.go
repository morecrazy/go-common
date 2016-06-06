package common

import (
	"fmt"
	"testing"
	"time"
)

type SqlCon struct {
	MysqlConn            string
	MysqlConnectPoolSize int
}

type FinanceServerConfigure struct {
	FinanceDb           SqlCon
	MallDb              SqlCon
	RedisNetwork        string
	RedisAddress        string
	RedisPassword       string
	RedisConnectTimeout int
	RedisReadTimeout    int
	RedisWriteTimeout   int
	RedisMaxActive      int
	RedisMaxIdle        int
	RedisIdleTimeout    int
	RedisWait           bool
	RedisDb             string
	GoodsServer         string
	RPCListen           string
	ServerMail          string
	FinanceAdminMails   string
	LogDir              string
	LogFile             string
}

func TestLoadCfg(t *testing.T) {
	start := time.Now()
	cfg := &FinanceServerConfigure{}
	addrs := []string{"http://120.26.17.34:4001"}
	err := LoadCfgFromEtcd(addrs, "financeserver", cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func(start time.Time) {
		fmt.Println(time.Now().Sub(start))
	}(start)
	fmt.Printf("%+v\n", cfg)
}

func TestLoadCfgContent(t *testing.T) {
	addrs := []string{"http://120.26.17.34:2379"}
	str, err := LoadContentFromEtcd(addrs, "qwebmiddleware", "/conf.d/ak")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", str)
}
