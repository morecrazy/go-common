package common

import . "backend/common/protocol"

var UserProfileClient *RpcClient
var UserRelationClient *RpcClient

func InitClient() error {
	var err error
	UserProfileClient, err = NewRpcClient(g_config.RpcSetting["UserProfileSetting"].Addr, g_config.RpcSetting["UserProfileSetting"].Net, UserprofileRpcFuncMap, "userprofile", Logger)
	if err != nil {
		Logger.Error("init UserProfileClient err :", err.Error())
	}

	UserRelationClient, err = NewRpcClient(g_config.RpcSetting["UserProfileSetting"].Addr, g_config.RpcSetting["UserProfileSetting"].Net, UserRelationRpcFuncMap, "userrelation", Logger)
	if err != nil {
		Logger.Error("init UserRelationClient err :", err.Error())
	}

	return nil
}
