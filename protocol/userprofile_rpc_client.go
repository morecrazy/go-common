package protocol

import "time"

type UserVipLabel struct {
	Id        int    `json:"id" gorm:"primary_key;column:id" sql:"auto_incremnet" pickle:"id"`
	Desc      string `gorm:"column:desc" json:"viplabel_desc";sql:type:"longtext" pickle:"viplabel_desc"`
	LIcon     string `gorm:"column:l_icon" json:"vipicon_l";sql:type:"varchar(128)" pickle:"vipicon_l"`
	IconPrior int    `gorm:"column:icon_prior" json:"vipicon_prior";sql:type:"smallint" pickle:"vipicon_prior"`
	MIcon     string `gorm:"column:m_icon" json:"vipicon_m";sql:type:"varchar(128)" pickle:"vipicon_m"`
	SIcon     string `gorm:"column:s_icon" json:"vipicon_s";sql:type:"varchar(128)" pickle:"vipicon_s"`
}

type UserProfile struct {
	Id               string            `json:"id" pickle:"id"`
	Nick             string            `json:"nick" pickle:"nick"`
	Realname         string            `json:"realname" pickle:"realname"`
	Gender           string            `json:"gender" pickle:"gender"`
	Location         string            `json:"location" pickle:"location"`
	Address          string            `json:"address" pickle:"address"`
	Portrait         string            `json:"portrait" pickle:"portrait"`
	Mobile_portraits []string          `json:"mobile_portraits" pickle:"mobile_portraits"`
	Birthday         map[string]int64  `json:"birthday" pickle:"birthday"`
	Descroption      string            `json:"descroption" pickle:"descroption"`
	Certificateinfo  string            `json:"certificateinfo" pickle:"certificateinfo"`
	Certificatename  string            `json:"certificatename" pickle:"certificatename"`
	Certificateid    string            `json:"certificateid" pickle:"certificateid"`
	Emailverified    bool              `json:"emailverified" pickle:"emailverified"`
	Mobileverified   bool              `json:"mobileverified" pickle:"mobileverified"`
	Followings       int64             `json:"followings" pickle:"followings"`
	Followers        int64             `json:"followers" pickle:"followers"`
	Tmp_portrait     string            `json:"tmp_portrait" pickle:"tmp_portrait"`
	Fighting_level   int64             `json:"fighting_level" pickle:"fighting_level"`
	Verify_code      string            `json:"verify_code" pickle:"verify_code"`
	Group_ids        string            `json:"group_ids" pickle:"group_ids"`
	Is_newuser       bool              `json:"is_newuser" pickle:"is_newuser"`
	Last_login       int64             `json:"last_login" pickle:"last_login"`
	Installed_apps   string            `json:"installed_apps" pickle:"installed_apps"`
	Hobby            string            `json:"hobby" pickle:"hobby"`
	Hobby_ids        string            `json:"hobby_ids" pickle:"hobby_ids"`
	Background_img   string            `json:"background_img" pickle:"background_img"`
	Unique_id        string            `json:"unique_id" pickle:"unique_id"`
	Mobilenumber     string            `json:"mobilenumber" pickle:"mobilenumber"`
	Domain           string            `json:"domain" pickle:"domain"`
	Email            string            `json:"email" pickle:"email"`
	Auto_id          int64             `json:"auto_id" pickle:"_auto_id"`
	Updated          int64             `json:"updated" pickle:"_updated"`
	T_Auto_id        int64             `json:"_auto_id" pickle:"auto_id"`
	T_Updated        int64             `json:"_updated" pickle:"updated`
	Is_vip           bool              `json:"is_vip" pickle:"is_vip"`
	PortraitMap      map[string]string `json:"portrait_map"`
	VipLabel         UserVipLabel      `json:"vip_label"`
	MobileUpdateTime time.Time         `json:mobileupdatetime`
}

type UserprofileList []UserProfile
type UserprofileBatchArgs struct {
	Ids []string `json:ids`
}

type UserprofileDefaultArgs struct {
	Id string `json:id`
}

type UserprofileDefaultReply struct {
	User UserProfile `json:"user"`
	//User userprofile.UserProfile `json:"user"`
}

type UserprofileBatchReply struct {
	//Users userprofile.UserprofileList `json:"users"`
	Users UserprofileList `json:"users"`
}

var UserprofileRpcFuncMap map[string]string = map[string]string{
	"batch_get": "UserprofileHandler.BatchGet",
	"get":       "UserprofileHandler.Get",
}
