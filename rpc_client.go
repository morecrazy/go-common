package common

import (
	. "backend/common/protocol"
	"fmt"
	"third/gorm"
)

var UserProfileClient *RpcClient
var UserLoginClient *RpcClient
var UserRelationClient *RpcClient
var RouteServerClinet *RpcClient
var SportsSortClinet *RpcClient

func InitRpcClient() error {
	var err error
	if nil == Config {
		return fmt.Errorf("please init common Conifg ")
	}
	UserProfileSetting, ok := Config.RpcSetting["UserProfileSetting"]
	if ok {
		UserProfileClient, err = NewRpcClient(UserProfileSetting.Addr, UserProfileSetting.Net, UserprofileRpcFuncMap, "userprofile", Logger)
		if err != nil {
			Logger.Error("init UserProfileClient err :", err.Error())
		}
		//userlogin 和 userprofile 同一个服务
		UserLoginClient, err = NewRpcClient(UserProfileSetting.Addr, UserProfileSetting.Net, UserloginRpcFuncMap, "userlogin", Logger)
		if err != nil {
			Logger.Error("init UserLoginClient err :", err.Error())
		}
	}

	UserRelationSetting, ok := Config.RpcSetting["UserRelationSetting"]
	if ok {
		UserRelationClient, err = NewRpcClient(UserRelationSetting.Addr, UserRelationSetting.Net, UserRelationRpcFuncMap, "userrelation", Logger)
		if err != nil {
			Logger.Error("init UserRelationClient err :", err.Error())
		}
	}

	RouteServerSetting, ok := Config.RpcSetting["RouteServerSetting"]
	if ok {
		RouteServerClinet, err = NewRpcClient(RouteServerSetting.Addr, RouteServerSetting.Net, RouteServerRpcFuncMap, "routeserver", Logger)
		if err != nil {
			Logger.Error("init RouteServerClinet err :", err.Error())
		}
	}

	SportsSortSetting, ok := Config.RpcSetting["SportsSortSetting"]
	if ok {
		SportsSortClinet, err = NewRpcClient(SportsSortSetting.Addr, SportsSortSetting.Net, SportsSortRpcFuncMap, "sportssort", Logger)
		if err != nil {
			Logger.Error("init SportsSortClinet err :", err.Error())
		}
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

func InitUserLoginClient(addr string, net string) error {
	var err error
	UserLoginClient, err = NewRpcClient(addr, net, UserloginRpcFuncMap, "userlogin", nil)
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
		if nil != Logger {
			Logger.Error(err.Error())
		}
		err = NewInternalError(RPCErrCode, err)
	}

	return reply.UserLogin, err
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
