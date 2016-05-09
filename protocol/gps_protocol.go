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
	Points        []interface{} `json:"points"`
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
	SportsType               int                    `json:"sports_type" form:"sports_type"`
	OffsetText               string                 `json:"offset_text" form:"offset_text"`
	SportsMode               int                    `json:"sportsMode" form:"sportsMode"`
	Version                  string                 `json:"version" form:"version"`
	HistoryVersion           int                    `json:"history_version" form:"history_version"`
	Location                 string                 `json:"location" form:"location"`
	StartTime                string                 `json:"start_time" form:"start_time"`
	EndTime                  string                 `json:"end_time" form:"end_time"`
	ProductId                string                 `json:"product_id" form:"product_id"`
	IsOpen                   int16                  `json:"is_open" form:"is_open"`
	IsReal                   interface{}            `json:"is_real" form:"is_real"` //ios定义为bool，android定义为int
	IsBaidu                  int16                  `json:"is_baidu" form:"is_baidu"`
	GoalValue                interface{}            `json:"goal_value" form:"goal_value"` //ios定义为int，android定义为float
	GoalType                 int                    `json:"goal_type" form:"goal_type"`
	StageDes                 string                 `json:"stage_des" form:"stage_des"`
	MaxAltitude              float64                `json:"MaxAltitude" form:"MaxAltitude"`
	MaxToPreviousSpeed       float64                `json:"MaxToPreviousSpeed" form:"MaxToPreviousSpeed"`
	TotalTime                interface{}            `json:"total_time" form:"total_time"` //ios定义为int，android定义为float
	TotalLength              float64                `json:"total_length" form:"total_length"`
	TotalCalories            float64                `json:"total_calories" form:"total_calories"`
	ActivityResult           int                    `json:"activity_result" form:"activity_result"`
	ActivityType             int                    `json:"activity_type" form:"activity_type"`
	BaiduCloud               interface{}            `json:"baidu_cloud" form:"baidu_cloud"` //ios定义为bool，android定义为int
	HighestSpeedPerkm        float64                `json:"highest_speed_perkm" form:"highest_speed_perkm"`
	CustomWords              string                 `json:"custom_words" form:"custom_words"`
	CaloriesPerm             []interface{}          `json:"calories_per_m" form:"calories_per_m"`
	ProgramName              string                 `json:"program_name" form:"program_name"`
	Model                    string                 `json:"model" form:"model"`
	LastOfProgram            int                    `json:"last_of_program" form:"last_of_program"`
	UserTimePerkm            []UseTimePerKm         `json:"usettime_per_km" form:"usettime_per_km"`
	UserStepsPerm            []interface{}          `json:"user_steps_list_perm" form:"user_steps_list_perm"`
	UserStepsValid           int16                  `json:"user_steps_valid" form:"user_steps_valid"`
	AverageSpeed             float64                `json:"AverageSpeed" form:"AverageSpeed"`
	Points                   []Point                `json:"points" form:"points"`
	HeartRate                map[string]interface{} `json:"heart_rate" form:"heart_rate"`
	GreenwayId               string                 `json:"greenway_id" form:"greenway_id"`
	HalfMarathon             interface{}            `json:"half_marathon" form:"half_marathon"`                             //ios定义为float，android定义为int
	Marathon                 interface{}            `json:"marathon" form:"marathon"`                                       //ios定义为float，android定义为int
	IsUserStopsportsAbnormal interface{}            `json:"is_user_stopsports_abnormal" form:"is_user_stopsports_abnormal"` //可能是android专用
	ReleaseVersion           interface{}            `json:"release_version" form:"release_version"`                         //可能是android专用
	StartDateTime            interface{}            `json:"StartDateTime" form:"StartDateTime"`                             //可能是android专用
	EndDateTime              interface{}            `json:"EndDateTime" form:"EndDateTime"`                                 //可能是android专用
	TotalTime1               interface{}            `json:"TotalTime" form:"TotalTime"`                                     //可能是android专用
	IsRoot                   interface{}            `json:"is_root" form:"is_root"`                                         //可能是android专用
	Locationcount            interface{}            `json:"locationcount" form:"locationcount"`                             //可能是android专用
	IsCrashRestore           interface{}            `json:"is_crash_restore" form:"is_crash_restore"`                       //可能是android专用
}

type Point struct {
	ToPreviousEnergy   float64 `json:"topreviousenergy" pickle:"topreviousenergy"`
	ToPreviousCostTime int     `json:"topreviouscostTime" pickle:"topreviouscostTime"`
	ToPreviousSpeed    float64 `json:"topreviousspeed" pickle:"topreviousspeed"`
	ToStartDistance    float64 `json:"tostartdistance" pickle:"tostartdistance"`
	ToStartCostTime    float64 `json:"tostartcostTime" pickle:"tostartcostTime"`
	Distance           float64 `json:"distance" pickle:"distance"`
	Longitude          float64 `json:"longitude" pickle:"longitude"`
	Latitude           float64 `json:"latitude" pickle:"latitude"`
	Elevation          float64 `json:"elevation" pickle:"elevation"`
	HAccuracy          float64 `json:"hAccuracy" pickle:"hAccuracy"`
	VAccuracy          float64 `json:"vAccuracy" pickle:"vAccuracy"`
	Type               int     `json:"type" pickle:"type"` //type 0: PROGRESS 1:SUSPEND else:RESTART
	TimeStamp          string  `json:"time_stamp" pickle:"time_stamp"`
	Steps              float64 `json:"-" pickle:"-"`
}

type UseTimePerKm struct {
	TotalUseTime float64    `json:"totalUseTime" pickle:"totalUseTime"`
	UseTime      float64    `json:"useTime" pickle:"useTime"`
	AtLocation   AtLocation `json:"atLocation" pickle:"atLocation"`
	Speed        float64    `json:"speed" pickle:"speed"`
	Distance     float64    `json:"distance" pickle:"distance"`
}

type AtLocation struct {
	Longitude float64 `json:"longitude" pickle:"longitude"`
	Latitude  float64 `json:"latitude" pickle:"latitude"`
}

type UserStepPerM struct {
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
