package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"third/context"
	etcd "third/etcd-client"
	//"time"
)

var RawData []byte

func NewEtcdApi(addrs []string) (etcd.KeysAPI, error) {
	cfg := etcd.Config{
		Endpoints: addrs,
		Transport: etcd.DefaultTransport,
		//HeaderTimeoutPerRequest: 10 * time.Second,
	}
	c, err := etcd.New(cfg)
	if err != nil {
		return nil, err
	}
	return etcd.NewKeysAPI(c), nil
}

func CfgFromEtcd(api etcd.KeysAPI, service, env string) (string, error) {
	rsp, err := api.Get(context.Background(), etcdKey(service, env), nil)
	if err != nil {
		log.Printf("read config [%s:%s] from etcd error:%v", service, env, err)
		return "", err
	}

	if rsp.Node == nil {
		log.Printf("empty etcd node")
		return "", errors.New("empty etcd node")
	}

	return rsp.Node.Value, nil
}

func LoadCfgFromEtcd(addrs []string, service,  cfg interface{}) error {
	api, err := NewEtcdApi(addrs)
	if err != nil {
		return err
	}

	data, err := CfgFromEtcd(api, service, "online")
	if err != nil {
		return err
	}
	RawData = []byte(data)

	return json.Unmarshal([]byte(data), cfg)
}

func etcdKey(service, env string) string {
	return fmt.Sprintf("/config/%s/%s", service, env)
}
