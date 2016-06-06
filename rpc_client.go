package common

import (
	"third/gorm"

	. "backend/common/protocol"
)

var UserProfileClient *RpcClient
var UserLoginClient *RpcClient
var UserRelationClient *RpcClient
var RouteServerClinet *RpcClient
var SportsSortClinet *RpcClient

func InitRpcClient() error {
	var err error
	UserProfileClient, err = NewRpcClient(Config.RpcSetting["UserProfileSetting"].Addr, Config.RpcSetting["UserProfileSetting"].Net, UserprofileRpcFuncMap, "userprofile", Logger)
	if err != nil {
		Logger.Error("init UserProfileClient err :", err.Error())
	}

	UserLoginClient, err = NewRpcClient(Config.RpcSetting["UserProfileSetting"].Addr, Config.RpcSetting["UserProfileSetting"].Net, UserloginRpcFuncMap, "userlogin", Logger)
	if err != nil {
		Logger.Error("init UserLoginClient err :", err.Error())
	}

	UserRelationClient, err = NewRpcClient(Config.RpcSetting["UserRelationSetting"].Addr, Config.RpcSetting["UserRelationSetting"].Net, UserRelationRpcFuncMap, "userrelation", Logger)
	if err != nil {
		Logger.Error("init UserRelationClient err :", err.Error())
	}

	RouteServerClinet, err = NewRpcClient(Config.RpcSetting["RouteServerSetting"].Addr, Config.RpcSetting["RouteServerSetting"].Net, RouteServerRpcFuncMap, "routeserver", Logger)
	if err != nil {
		Logger.Error("init RouteServerClinet err :", err.Error())
	}

	SportsSortClinet, err = NewRpcClient(Config.RpcSetting["SportsSortSetting"].Addr, Config.RpcSetting["SportsSortSetting"].Net, SportsSortRpcFuncMap, "sportssort", Logger)
	if err != nil {
		Logger.Error("init SportsSortClinet err :", err.Error())
	}

	return nil
}

// 只使用 GetProfile的接口，单独初始化
func InitProfileClient(addr string, net string) error {
	var err error
	UserProfileClient, err = NewRpcClient(addr, net, UserprofileRpcFuncMap, "userprofile", nil)
	if err != nil {
		return err
	}

	return nil
}

func GetProfileById(userId string) (UserProfile, error) {
	if userId == "" {
		return UserProfile{}, gorm.RecordNotFound
	}
	var reply UserprofileDefaultReply
	args := UserprofileDefaultArgs{
		Id: userId,
	}
	//	Logger.Debug("GetProfileById %v", args)
	err := UserProfileClient.Call("get", &args, &reply)
	if err != nil {
		if nil != Logger {
			Logger.Error(err.Error())
		}
		err = NewInternalError(RPCErrCode, err)
	}

	return reply.User, err
}

func BatchGetProfileByIds(userIds []string) (UserprofileList, error) {
	var reply UserprofileBatchReply
	args := UserprofileBatchArgs{
		Ids: userIds,
	}
	err := UserProfileClient.Call("batch_get", &args, &reply)
	if err != nil {
		if nil != Logger {
			Logger.Error(err.Error())
		}
		err = NewInternalError(RPCErrCode, err)
	}
	return reply.Users, err
}

func GetLoginById(userId string) (UserLogin, error) {
	if userId == "" {
		return UserLogin{}, gorm.RecordNotFound
	}
	var reply UserloginDefaultReply
	args := UserloginDefaultArgs{
		Id: userId,
	}
	//	Logger.Debug("GetProfileById %v", args)
	err := UserLoginClient.Call("get", &args, &reply)
	if err != nil {
		Logger.Error(err.Error())
		err = NewInternalError(RPCErrCode, err)
	}

	return reply.UserLogin, err
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

func GetFollowing(user_id string) (following_ids []string, err error) {
	var get_all = false
	var page_num = 1
	for !get_all {
		var req = GetFollowingReq{
			Selfuserid: user_id,
			Pagenum:    page_num,
			Pagesize:   1000,
		}
		var resp GetFollowingRes

		err := UserRelationClient.Call("get_following", &req, &resp)
		if nil != err {
			Logger.Error("get_following err :%v", err)
			return following_ids, err
		}
		page_num += 1

		if 0 == len(resp.Data) {
			break
		}

		for _, value := range resp.Data {
			following_ids = append(following_ids, value)
		}
		if len(resp.Data) >= 1000 {
			get_all = false
		} else {
			get_all = true
		}
	}

	return following_ids, err
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

func ProcRouteLog(routeId, userId string, postData map[string]interface{}) (SaveRouteLogRes, error) {
	var reply SaveRouteLogRes
	args := SaveRouteLogReq{
		RouteLog: RouteLog{
			RouteId:  routeId,
			UserId:   userId,
			PostData: postData,
		},
	}
	//	Logger.Debug("ProcRouteLog arg %v", args)
	err := RouteServerClinet.Call("process_routelog", &args, &reply)
	if err != nil {
		Logger.Error(err.Error())
		err = NewInternalError(RPCErrCode, err)
	}

	return reply, err
}

func DeleteRoute(routeId, userId string) (DeleteRouteRes, error) {
	var reply DeleteRouteRes
	args := DeleteRouteReq{
		RouteId: routeId,
		UserId:  userId,
	}
	//	Logger.Debug("ProcRouteLog arg %v", args)
	err := RouteServerClinet.Call("delete_route", &args, &reply)
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
	//	Logger.Debug("UpdateSportInfo arg %v", args)
	err := SportsSortClinet.Call("update_sport_info", &args, &reply)
	if err != nil {
		Logger.Error(err.Error())
		err = NewInternalError(RPCErrCode, err)
	}

	return reply, err
}

func GetUserSort(userIds []string, sportType, sortType, relationType, pageNum, pageSize int) (GetUserSortResp, error) {
	var reply GetUserSortResp
	args := GetUserSortReq{
		UserIds:      userIds,
		SportType:    sportType,
		SortType:     sortType,
		RelationType: relationType,
		PageNum:      pageNum,
		PageSize:     pageSize,
	}
	//	Logger.Debug("GetUserSort arg %v", args)
	err := SportsSortClinet.Call("get_user_sport", &args, &reply)
	if err != nil {
		Logger.Error(err.Error())
		err = NewInternalError(RPCErrCode, err)
	}

	return reply, err
}
