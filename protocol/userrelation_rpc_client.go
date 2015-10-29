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

var UserRelationRpcFuncMap map[string]string = map[string]string{
	"get_follower":        "UserRelationHandler.GetFollower",
	"get_following_flag":  "UserRelationHandler.GetFollowingFlag",
	"get_following_count": "UserRelationHandler.GetFollowedCount",
	"get_following_flags": "UserRelationHandler.GetFollowingFlags",
	"get_following":       "UserRelationHandler.GetFollowing",
}
