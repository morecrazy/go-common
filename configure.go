package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type MysqlConfig struct {
	MysqlConn            string
	MysqlConnectPoolSize int
}

type RedisConfig struct {
	RedisConn      string
	RedisPasswd    string
	ReadTimeout    int
	ConnectTimeout int
	WriteTimeout   int
	IdleTimeout    int
	MaxIdle        int
	MaxActive      int
	RedisDb        string
}

type RPCSetting struct {
	Addr string
	Net  string
}

type CeleryQueue struct {
	Url   string
	Queue string
}

type OssConfig struct {
	AccessKeyId     string
	AccessKeySecret string
	Region          string
	Bucket          string
}

type Configure struct {
	MysqlSetting  map[string]MysqlConfig
	RedisSetting  map[string]RedisConfig
	RpcSetting    map[string]RPCSetting
	CelerySetting map[string]CeleryQueue
	OssSetting    map[string]OssConfig
	SentryUrl     string
	LogDir        string
	LogFile       string
	Listen        string
	RpcListen     string
	LogLevel      string
	External      map[string]string
	ExternalInt64 map[string]int64
}

var Config *Configure
var g_config_file_last_modify_time time.Time
var g_local_conf_file string

/*
func (config Configure) String() string {
	return fmt.Sprintf("[Mysqlconn:%s][RedisTemp:%s][RedisSta:%s][RedisBig:%s][RedisTimer:%s][LogFile:%s][MysqlConnectPoolSize:%d] [UserRelationSetting:%s][UserProfileSetting:%s]", config.MysqlFeed.MysqlConn, config.RedisTemp.RedisConn, config.RedisSta.RedisConn, config.RedisBig.RedisConn, config.RedisTimer.RedisConn, config.LogFile, config.MysqlFeed.MysqlConnectPoolSize, config.UserRelationSetting.Addr, config.UserProfileSetting.Addr)
}
*/
func InitConfigFile(filename string, config *Configure) error {
	//var environment = os.Getenv("GOENV")
	//
	//fmt.Println("environment", environment)
	//switch environment {
	//case "DEV":
	//	filename = filename + ".dev"
	//case "TEST":
	//	filename = filename + ".test"
	//case "PRE":
	//	filename = filename + ".pre"
	//case "ONLINE":
	//	filename = filename + ".online"
	//default:
	//	filename = filename + ".test"
	//}

	fmt.Println("filename", filename)
	fi, err := os.Stat(filename)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return err
	}

	if g_config_file_last_modify_time.Equal(fi.ModTime()) {
		return nil //fmt.Errorf("Config File Have No Change")
	} else {
		g_config_file_last_modify_time = fi.ModTime()
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return err
	}

	if err := json.Unmarshal(bytes, config); err != nil {
		err = NewInternalError(DecodeErrCode, err)
		fmt.Println("Unmarshal: ", err.Error())
		return err
	}
	fmt.Println("conifg :", *config)
	g_local_conf_file = filename
	Config = config
	return nil
}

func RunDynamicConfigureTimer() {
	confTimer := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-confTimer.C:

			dynamicConfigure()
		}
	}
}

func dynamicConfigure() {
	err := InitConfigFile(g_local_conf_file, Config)
	if err != nil {
		fmt.Println(err)
		return
	}

	ChangeLogLevel(Config.LogLevel)
}
