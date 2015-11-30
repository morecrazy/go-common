package common

import . "backend/common/protocol"

var UserProfileClient *RpcClient
var UserRelationClient *RpcClient
var RouteServerClinet *RpcClient

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

	RouteServerClinet, err = NewRpcClient(Config.RpcSetting["RouteServerSetting"].Addr, Config.RpcSetting["RouteServerSetting"].Net, RouteServerRpcFuncMap, "routeserver", Logger)
	if err != nil {
		Logger.Error("init RouteServerClinet err :", err.Error())
	}

	return nil
}

func SimplifyProcRouteLog(userId string, postData map[string]interface{}) (SimplifyProcRouteLogRes, error) {
	var reply SimplifyProcRouteLogRes
	args := SimplifyProcRouteLogReq{
		RouteLog: RouteLog{
			UserId:   userId,
			PostData: postData,
		},
	}
	//	Logger.Debug("SimplifyProcRouteLog arg %v", args)
	err := RouteServerClinet.Call("simplify_proc_route_log", &args, &reply)
	if err != nil {
		Logger.Error(err.Error())
		err = NewInternalError(RPCErrCode, err)
	}

	return reply, err
}

func SaveRouteLog(routeId, userId string, postData map[string]interface{}) (SaveRouteLogRes, error) {
	var reply SaveRouteLogRes
	args := SaveRouteLogReq{
		RouteLog: RouteLog{
			RouteId:  routeId,
			UserId:   userId,
			PostData: postData,
		},
	}
	//	Logger.Debug("SimplifyProcRouteLog arg %v", args)
	err := RouteServerClinet.Call("save_routelog", &args, &reply)
	if err != nil {
		Logger.Error(err.Error())
		err = NewInternalError(RPCErrCode, err)
	}

	return reply, err
}
