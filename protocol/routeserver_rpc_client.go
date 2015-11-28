package protocol

type SimplifyProcRouteLogReq struct {
	RouteLog RouteLog `json:"route_log"`
}

type SimplifyProcRouteLogRes struct {
	UserId  string `json:"user_id"`
	Id      int64  `json:"id"`
	IsFraud bool   `json:"is_fraud"`
}

var UserRelationRpcFuncMap map[string]string = map[string]string{
	"simplify_proc_route_log": "RouteServerClinet.SimplifyProcRouteLog",
}
