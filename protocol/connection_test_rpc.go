package protocol

type PingReq struct {
	Id      int64  `json:"id"`
	Src     string `json:"src"`
	Dst     string `json:"dst"`
	ReqTime string `json:"req_time"`
	Len     int    `json:"len"`
	PayLoad string `json:"pay_load"`
}

type PingResp struct {
	Id       int64  `json:"id"`
	Src      string `json:"src"`
	Dst      string `json:"dst"`
	RealDst  string `json:"real_dst"`
	Len      int    `json:"len"`
	RespTime string `json:"resp_time"`
	PayLoad  string `json:"pay_load"`
}

var ConnectionTestRpcFuncMap map[string]string = map[string]string{
	"ping": "ConnectionTestHandler.Ping",
}
