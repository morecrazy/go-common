package common

import (
	"backend/jsonrpc"
	"backend/rpc"
	"net"
	"time"
)

type Log interface {
	Error(format string, args ...interface{})
	Info(format string, args ...interface{})
	Notice(format string, args ...interface{})
	Debug(format string, args ...interface{})
}

type RpcClient struct {
	rpc_client *rpc.Client
	Addr       string
	Net        string
	name       string
	func_map   map[string]string
	logger     Log
	pool       *rpc.Pool
}

func NewRpcClient(addr, net string, func_map map[string]string, name string,
	logger Log) (*RpcClient, error) {
	var err error
	client := &RpcClient{}
	client.Addr = addr
	client.Net = net
	client.func_map = func_map
	client.name = name
	client.logger = logger
	client.pool = rpc.NewPool(client.Connect, 100, 100, 1000*time.Second, true)

	return client, err
}

func (client *RpcClient) Connect() (*rpc.Client, error) {
	conn, err := net.DialTimeout(client.Net, client.Addr, 2*time.Second)
	if err != nil {
		client.logger.Error("get %s rpc client error :%v", client.name, err)
		return nil, err
	}

	rpc_client := jsonrpc.NewClient(conn)

	return rpc_client, nil
}

func (client *RpcClient) Call(method string, args interface{}, reply interface{}) error {
	if "add_picture" != method {
		if nil != client.logger {
			client.logger.Debug("call rpc : %s, %v", method, args)
		}
	}
	rpc_client, err := client.pool.Get()
	if err != nil {
		if nil != client.logger {
			client.logger.Error("get %s rpc client error :%v", client.name, err)
		}
		return err
	}
	defer rpc_client.Close()
	err = rpc_client.CallTimeout(client.func_map[method], args, reply, 2, 2*time.Second)
	if err != nil {
		is_user_err, _, _ := IsUserErr(err)
		if is_user_err {
			if nil != client.logger {
				client.logger.Notice("call %s rpc client err :%s,%v", client.name, method, err)
			}
		} else {
			if nil != client.logger {
				client.logger.Error("call %s rpc client err :%s,%v", client.name, method, err)
			}
		}
	}

	return err
}

func (client *RpcClient) DirectCall(method string, args interface{}, reply interface{}) error {
	if "add_picture" != method {
		client.logger.Info("call rpc : %s, %v", method, args)
	}

	err := client.rpc_client.CallTimeout(client.func_map[method], args, reply, 2*time.Second)
	if err != nil {
		is_user_err, _, _ := IsUserErr(err)
		if is_user_err {
			client.logger.Notice("call %s rpc client err :%s,%v", client.name, method, err)
		} else {
			client.logger.Error("call %s rpc client err :%s,%v", client.name, method, err)
		}
	}

	return err
}
