package protocol

type SimplifyProcRouteLogReq struct {
	RouteLog RouteLog `json:"route_log"`
}

type SimplifyProcRouteLogRes struct {
	UserId  string `json:"user_id"`
	Id      int64  `json:"id"`
	IsFraud bool   `json:"is_fraud"`
}

type SaveRouteLogReq struct {
	RouteLog RouteLog `json:"route_log"`
}

type SaveRouteLogRes struct {
	RouteId string `json:"route_id"`
	UserId  string `json:"user_id"`
	Id      int64  `json:"id"`
	IsFraud bool   `json:"is_fraud"`
}

type DeleteRouteReq struct {
	UserId  string `json:"user_id"`
	RouteId string `json:"route_id"`
}

type DeleteRouteRes struct {
	Status bool `json:"status"`
}

var RouteServerRpcFuncMap map[string]string = map[string]string{
	"simplify_proc_route_log": "RouteHandler.SimplifyProcRouteLog",
	"process_routelog":        "RouteHandler.ProcessRouteLog",
	"delete_route":            "RouteHandler.DeleteRoute",
}
