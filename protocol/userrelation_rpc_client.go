package protocol

type GetFollowingReq struct {
	Selfuserid string
	Pagenum    int
	Pagesize   int
}

type GetFollowingRes struct {
	Data []string
}

type GetFlagRes struct {
	Data int
}

type FollowingReq struct {
	Selfuserid   string
	Targetuserid string
}

type FollowingRes struct {
	Data bool
}

type GetCountRes struct {
	Data int
}

type GetFlagsRes struct {
	Data map[string]int
}

var UserRelationRpcFuncMap map[string]string = map[string]string{
	"follow":              "UserRelationHandler.FollowPeople",
	"unfollow":            "UserRelationHandler.UnFollowPeople",
	"get_follower":        "UserRelationHandler.GetFollower",
	"get_following_flag":  "UserRelationHandler.GetFollowingFlag",
	"get_following_count": "UserRelationHandler.GetFollowingCount",
	"get_follower_count":  "UserRelationHandler.GetFollowedCount",
	"get_friends_count":   "UserRelationHandler.GetFriendsCount",
	"get_following_flags": "UserRelationHandler.GetFollowingFlags",
	"get_following":       "UserRelationHandler.GetFollowing",
	"get_friends":         "UserRelationHandler.GetFriends",
}
