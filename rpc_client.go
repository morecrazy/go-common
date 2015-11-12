package common

import . "backend/common/protocol"

var UserProfileClient *RpcClient
var UserRelationClient *RpcClient

func InitClient() error {
	var err error
	UserProfileClient, err = NewRpcClient(Config.RpcSetting["UserProfileSetting"].Addr, Config.RpcSetting["UserProfileSetting"].Net, UserprofileRpcFuncMap, "userprofile", Logger)
	if err != nil {
		Logger.Error("init UserProfileClient err :", err.Error())
	}

	UserRelationClient, err = NewRpcClient(Config.RpcSetting["UserProfileSetting"].Addr, Config.RpcSetting["UserProfileSetting"].Net, UserRelationRpcFuncMap, "userrelation", Logger)
	if err != nil {
		Logger.Error("init UserRelationClient err :", err.Error())
	}

	return nil
}
