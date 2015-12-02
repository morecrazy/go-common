package protocol

type UpdateUserSportInfoReq struct {
	UserId       string  `json:"user_id"`
	CurDay       string  `json:"cur_day"`
	DaySummary   float64 `json:"day_summary"`
	WeekSummary  float64 `json:"week_summary"`
	MonthSummary float64 `json:"month_summary"`
	YearSummary  float64 `json:"year_summary"`
	AllSummary   float64 `json:"all_summary"`
}

type UpdateUserSportInfoResp struct {
	Data int `json:"data"`
}

var SportSortRpcFuncMap map[string]string = map[string]string{
	"update_sport_info": "SportsSortHandler.UpdateUserSportInfo",
	//"save_routelog":           "RouteHandler.SaveRouteLog",
}
