package protocol

import "time"

type Content struct {
	Step      int64   `json:"step" pickle:"steps"`
	Calories  float64 `json:"calories" pickle:"calories"`
	Distance  float64 `json:"distance" pickle:"distance"`
	State     bool    `json:"state" pickle:"state"` //true:无效 false:有效
	TotalTime float64 `json:"total_time" pickle:"total_time"`
}

type HeartRate struct {
	Id      int64                  `json:"id"`
	RouteId string                 `json:"route_id"`
	UserId  string                 `json:"user_id"`
	MaxRate int32                  `json:"max_rate"`
	AvgRate int32                  `json:"avg_rate"`
	Content map[string]interface{} `json:"content"`
}

type RouteData struct {
	Id            int64         `json:"id"`
	RouteId       string        `json:"route_id"`
	UserId        string        `json:"user_id"`
	UseTtimePerKm []interface{} `json:"usettime_per_km"`
	RouteLine     []interface{} `json:"route_line"`
	Points        []interface{} `json"points"`
}

type RouteDayDetail struct {
	UserId      string          `json:"user_id"`
	SportsType  int             `json:"sports_type"`
	DayContent  map[int]Content `json:"content"`
	ProductId   string          `json:"product_id"`
	Curday      time.Time       `json:"curday"`
	TotalTime   float64         `json:"total_time"`
	TotalLength float64         `json:"total_length"`
	TotalCalory float64         `json:"total_calory"`
	TotalStep   int64           `json:"total_step"`
}

type RouteLog struct {
	Id       int64
	RouteId  string                 `json:"route_id"`
	UserId   string                 `json:"user_id"`
	PostData map[string]interface{} `json:"post_data"`
}

type PostData struct {
	StartTime     string      `json:"start_time"`
	EndTime       string      `json:"end_time"`
	ProductId     string      `json:"product_id"`
	SportsType    float64     `json:"sports_type"`
	TotalLength   float64     `json:"total_length"`
	TotalCalories float64     `json:"total_calories"`
	TotalTime     float64     `json:"total_time"`
	AverageSpeed  float64     `json:"AverageSpeed"`
	Points        interface{} `json:"points"`
	UserTimePerkm interface{} `json:"usettime_per_km"`
	UserSteps     interface{} `json:"user_steps_list"`
	UserStepsPerm interface{} `json:"user_steps_list_perm"`
}

type Point struct {
	ToPreviousEnergy   int     `json:topreviousenergy"`
	ToPreviousCostTime int     `json:"topreviouscostTime"`
	ToPreviousSpeed    float64 `json:"topreviousspeed"`
	ToStartDistance    int     `json:"tostartdistance"`
	ToStartCostTime    float64 `json:"tostartcostTime"`
	Distance           float64 `json:"distance"`
	Longitude          float64 `json:"longitude"`
	Latitude           float64 `json:"latitude"`
	Elevation          float64 `json:"elevation"`
	HAccuracy          float64 `json:"hAccuracy"`
	VAccuracy          int     `json:"vAccuracy"`
	Type               int     `json:"type"` //type 0: PROGRESS 1:SUSPEND else:RESTART
	TimeStamp          string  `json:"time_stamp"`
}

type UseTimePerKm struct {
	TotalUseTime float64 `json:"totalUseTime"`
	UseTime      float64 `json:"useTime"`
	Speed        float64 `json:"speed"`
	Distance     float64 `json:"distance"`
}

type UserStep struct {
	Timestamp string `json:"time_stamp"`
	Steps     int64  `json:"steps"`
}

type MessageData struct {
	UserId      string    `json:"user_id"`
	RouteId     string    `json:"route_id"`
	UpLoadTime  time.Time `json:"upload_time"`
	TotalTime   float64   `json:"total_time"`
	TotalLength float64   `json:"total_length"`
	TotalCalory float64   `json:"total_calory"`
	IsFraud     bool      `json:"is_fraud"`
}

type SendData struct {
	MessageType string      `json:"message_type"`
	Data        MessageData `json:"data"`
}

type UserSumData struct {
	TotalTime   float64
	TotalLength float64
	TotalCalory float64
	TotalStep   int64
}

type StepContent struct {
	Time     string  `json:time`
	Steps    int64   `json:steps`
	Calories float64 `json:calories`
	Meters   float64 `json:meters`
}

type MergeContent struct {
	Step      int64   `json:"steps" pickle:"steps"`
	Calories  float64 `json:"calories" pickle:"calories"`
	Distance  float64 `json:"distance" pickle:"distance"`
	TotalTime float64 `json:"total_time" pickle:"total_time"`
}
