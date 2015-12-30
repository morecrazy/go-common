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

type GetUserSortReq struct {
	UserIds      []string `json:"user_ids"`
	SportType    int      `json:"sport_type"`    //0 all, 1 run, 2 ride, 3 walk, 4 ski, 5 skate, 6 climb
	SortType     int      `json:"sort_type"`     //0 week, 1 month, 2 all, 3 last week  4, day
	RelationType int      `json:"relation_type"` //0 group, 1 friend,
	PageNum      int      `json:"page_num"`
	PageSize     int      `json:"page_size"`
}

type UserSportInfo struct {
	UserId   string  `json:"user_id"`
	Distance float64 `json:"distance"`
}

type SportSlice []*UserSportInfo

type GetUserSortResp struct {
	Data SportSlice `json:"data"`
}

var SportsSortRpcFuncMap map[string]string = map[string]string{
	"update_sport_info": "SportsSortHandler.UpdateUserSportInfo",
	"get_user_sort":     "SportsSortHandler.GetUserSort",
}
