package common

import . "backend/common/protocol"

var UserProfileClient *RpcClient
var UserRelationClient *RpcClient
var RouteServerClinet *RpcClient
var SportSortClinet *RpcClient

func InitClient() error {
	var err error
	UserProfileClient, err = NewRpcClient(Config.RpcSetting["UserProfileSetting"].Addr, Config.RpcSetting["UserProfileSetting"].Net, UserprofileRpcFuncMap, "userprofile", Logger)
	if err != nil {
		Logger.Error("init UserProfileClient err :", err.Error())
	}

	UserRelationClient, err = NewRpcClient(Config.RpcSetting["UserRelationSetting"].Addr, Config.RpcSetting["UserRelationSetting"].Net, UserRelationRpcFuncMap, "userrelation", Logger)
	if err != nil {
		Logger.Error("init UserRelationClient err :", err.Error())
	}

	RouteServerClinet, err = NewRpcClient(Config.RpcSetting["RouteServerSetting"].Addr, Config.RpcSetting["RouteServerSetting"].Net, RouteServerRpcFuncMap, "routeserver", Logger)
	if err != nil {
		Logger.Error("init RouteServerClinet err :", err.Error())
	}

	SportSortClinet, err = NewRpcClient(Config.RpcSetting["SportSortSetting"].Addr, Config.RpcSetting["SportSortSetting"].Net, SportSortRpcFuncMap, "sportsort", Logger)
	if err != nil {
		Logger.Error("init SportSortClinet err :", err.Error())
	}

	return nil
}

func GetProfileById(userId string) (UserProfile, error) {
	var reply UserprofileDefaultReply
	args := UserprofileDefaultArgs{
		Id: userId,
	}
	Logger.Debug("GetProfileById %v", args)
	err := UserProfileClient.Call("get", &args, &reply)
	if err != nil {
		Logger.Error(err.Error())
		err = NewInternalError(RPCErrCode, err)
	}

	return reply.User, err
}

func GetFollower(userId string) (follower_ids []string, err error) {
	var get_all = false
	var page_num = 1
	for !get_all {
		var req = GetFollowingReq{
			Selfuserid: userId,
			Pagenum:    page_num,
			Pagesize:   1000,
		}
		var resp GetFollowingRes

		err = UserRelationClient.Call("get_follower", &req, &resp)
		if nil != err {
			Logger.Error("get_follower err :%v", err)
			return follower_ids, err
		}
		page_num += 1

		if 0 == len(resp.Data) {
			break
		}

		for _, value := range resp.Data {
			follower_ids = append(follower_ids, value)
		}
		if len(resp.Data) >= 1000 {
			get_all = false
		} else {
			get_all = true
		}
	}

	return follower_ids, err
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

func UpdateSportInfo(userId, curDay string, daySummary, weekSummary, monthSummary, yearSummary, allSummary float64) (UpdateUserSportInfoResp, error) {
	var reply UpdateUserSportInfoResp
	args := UpdateUserSportInfoReq{
		UserId:       userId,
		CurDay:       curDay,
		DaySummary:   daySummary,
		WeekSummary:  weekSummary,
		MonthSummary: monthSummary,
		YearSummary:  yearSummary,
		AllSummary:   allSummary,
	}
	Logger.Debug("UpdateSportInfo arg %v", args)
	err := SportSortClinet.Call("update_sport_info", &args, &reply)
	if err != nil {
		Logger.Error(err.Error())
		err = NewInternalError(RPCErrCode, err)
	}

	return reply, err
}
