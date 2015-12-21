package protocol

type UserLogin struct {
	AutoId         int64  `gorm:"primary_key;column:auto_id" sql:"auto_increment"`
	Id             string `gorm:"column:id" sql:"type:varchar(36);unique:id" json:"id"`
	Email          string `gorm:"column:email" sql:"type:varchar(50);unique:user_login_idx_email2" json:"email"`
	Mobilenumber   string `gorm:"column:mobilenumber" sql:"type:varchar(15);unique:user_loginI3" json:"mobilenumber"`
	Password       string `gorm:"column:password" sql:"type:varchar(51)" json:"password"`
	CanLogin       int    `gorm:"column:canlogin" sql:"type:tinyint(1)" json:"canlogin"`
	TimeZone       int    `gorm:"column:time_zone" sql:"type:int(11)" json:"time_zone"`
	RegistDatetime int    `gorm:"column:regist_datetime" sql:"type:int(11)" json:"regist_datetime"`
	LastLogin      int    `gorm:"column:last_login" sql:"type:int(11)" json:"last_login"`
}

type UserloginDefaultArgs struct {
	Id string `json:"id"`
}

type UserloginDefaultReply struct {
	UserLogin UserLogin `json:"userlogin"`
}

var UserloginRpcFuncMap map[string]string = map[string]string{
	"get": "UserLoginHandler.Get",
}
