package common

import (
	"backend/jsonrpc"
	"backend/rpc"
	"net"
)

type Log interface {
	Error(format string, args ...interface{})
	Info(format string, args ...interface{})
	Notice(format string, args ...interface{})
}

type RpcClient struct {
	rpc_client *rpc.Client
	Addr       string
	Net        string
	name       string
	func_map   map[string]string
	logger     Log
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

	return client, err
}

func (client *RpcClient) Connect(net_type, addr string) (*rpc.Client, error) {
	conn, err := net.Dial(net_type, addr)
	if err != nil {
		client.logger.Error("get %s rpc client error :%v", client.name, err)
		return nil, err
	}

	rpc_client := jsonrpc.NewClient(conn)
	return rpc_client, nil
}

func (client *RpcClient) Call(method string, args interface{}, reply interface{}) error {
	if "add_picture" != method {
		client.logger.Info("call rpc : %s, %v", method, args)
	}
	conn, err := net.Dial(client.Net, client.Addr)
	if err != nil {
		client.logger.Error("get %s rpc client error :%v", client.name, err)
		return err
	}
	defer conn.Close()
	rpc_client := jsonrpc.NewClient(conn)
	defer rpc_client.Close()
	err = rpc_client.Call(client.func_map[method], args, reply)
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
