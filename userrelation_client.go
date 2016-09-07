package common

import (
	. "backend/common/protocol"
	"fmt"
)

func InitRelationRpcClient(addr string, net string) error {
	var err error
	UserRelationClient, err = NewRpcClient(addr, net, UserRelationRpcFuncMap, "userrelation", Logger)
	if err != nil {
		Logger.Error("init UserRelationClient err :", err.Error())
	}
	return err
}

//需要验证userid，target_user_id有效性
func FollowPeople(user_id, target_user_id string) error {
	if "" == user_id || "" == target_user_id || user_id == target_user_id {
		return fmt.Errorf("%s", "params error")
	}
	var followReq FollowingReq
	var followRes FollowingRes
	followReq.Selfuserid = user_id
	followReq.Targetuserid = target_user_id

	err := UserRelationClient.Call("follow", &followReq, &followRes)
	if err != nil {
		Logger.Error(err.Error())
		return err
	}

	return nil
}
//需要验证userid，target_user_id有效性
func UnFollowPeople(user_id, target_user_id string) error {
	if "" == user_id || "" == target_user_id || user_id == target_user_id {
		return fmt.Errorf("%s", "params error")
	}
	var followReq FollowingReq
	var followRes FollowingRes
	followReq.Selfuserid = user_id
	followReq.Targetuserid = target_user_id

	err := UserRelationClient.Call("unfollow", &followReq, &followRes)
	if err != nil {
		Logger.Error(err.Error())
		return err
	}

	return nil
}

func GetFollower(userId string) (follower_ids []string, err error) {
	var get_all = false
	var page_num = 1
	for !get_all {
		var req = GetFollowingReq{
			Selfuserid: userId,
			Pagenum:    page_num,
			Pagesize:   1000,
		}
		var resp GetFollowingRes

		err = UserRelationClient.Call("get_follower", &req, &resp)
		if nil != err {
			Logger.Error("get_follower err :%v", err)
			return follower_ids, err
		}
		page_num += 1

		if 0 == len(resp.Data) {
			break
		}

		for _, value := range resp.Data {
			follower_ids = append(follower_ids, value)
		}
		if len(resp.Data) >= 1000 {
			get_all = false
		} else {
			get_all = true
		}
	}

	return follower_ids, err
}

func GetFollowerByPage(userId string, page_num, limit int) (follower_ids []string, err error) {

	if limit > 1000 {
		return follower_ids, fmt.Errorf("can't get follower more than 1000 once time")
	}

	var req = GetFollowingReq{
		Selfuserid: userId,
		Pagenum:    page_num,
		Pagesize:   limit,
	}
	var resp GetFollowingRes

	err = UserRelationClient.Call("get_follower", &req, &resp)
	if nil != err {
		Logger.Error("get_follower err :%v", err)
		return follower_ids, err
	}

	if 0 == len(resp.Data) {
		return follower_ids, err
	}

	follower_ids = make([]string, 0, len(resp.Data))
	for _, value := range resp.Data {
		follower_ids = append(follower_ids, value)
	}

	return follower_ids, err
}

func GetFollowing(user_id string) (following_ids []string, err error) {
	var get_all = false
	var page_num = 1
	for !get_all {
		var req = GetFollowingReq{
			Selfuserid: user_id,
			Pagenum:    page_num,
			Pagesize:   1000,
		}
		var resp GetFollowingRes

		err := UserRelationClient.Call("get_following", &req, &resp)
		if nil != err {
			Logger.Error("get_following err :%v", err)
			return following_ids, err
		}
		page_num += 1

		if 0 == len(resp.Data) {
			break
		}

		for _, value := range resp.Data {
			following_ids = append(following_ids, value)
		}
		if len(resp.Data) >= 1000 {
			get_all = false
		} else {
			get_all = true
		}
	}

	return following_ids, err
}

func GetFollowingByPage(user_id string, page_num, limit int) (following_ids []string, err error) {
	if limit > 1000 {
		return following_ids, fmt.Errorf("can't get following more than 1000 once time")
	}

	var req = GetFollowingReq{
		Selfuserid: user_id,
		Pagenum:    page_num,
		Pagesize:   limit,
	}
	var resp GetFollowingRes

	err = UserRelationClient.Call("get_following", &req, &resp)
	if nil != err {
		Logger.Error("get_following err :%v", err)
		return following_ids, err
	}
	if 0 == len(resp.Data) {
		return following_ids, err
	}

	following_ids = make([]string, 0, len(resp.Data))
	for _, value := range resp.Data {
		following_ids = append(following_ids, value)
	}

	return following_ids, err
}

func GetFriends(user_id string) (friends_ids []string, err error) {
	var get_all = false
	var page_num = 1
	for !get_all {
		var req = GetFollowingReq{
			Selfuserid: user_id,
			Pagenum:    page_num,
			Pagesize:   1000,
		}
		var resp GetFollowingRes

		err := UserRelationClient.Call("get_friends", &req, &resp)
		if nil != err {
			Logger.Error("get_friends err :%v", err)
			return friends_ids, err
		}
		page_num += 1

		if 0 == len(resp.Data) {
			break
		}

		for _, value := range resp.Data {
			friends_ids = append(friends_ids, value)
		}
		if len(resp.Data) >= 1000 {
			get_all = false
		} else {
			get_all = true
		}
	}

	return friends_ids, err
}

func GetFriendsByPage(user_id string, page_num, limit int) (friends_ids []string, err error) {
	if limit > 1000 {
		return friends_ids, fmt.Errorf("can't get friends more than 1000 once time")
	}
	var req = GetFollowingReq{
		Selfuserid: user_id,
		Pagenum:    page_num,
		Pagesize:   1000,
	}
	var resp GetFollowingRes

	err = UserRelationClient.Call("get_friends", &req, &resp)
	if nil != err {
		Logger.Error("get_friends err :%v", err)
		return friends_ids, err
	}

	if 0 == len(resp.Data) {
		return friends_ids, err
	}

	friends_ids = make([]string, 0, len(resp.Data))
	for _, value := range resp.Data {
		friends_ids = append(friends_ids, value)
	}

	return friends_ids, err
}

// 返回值 0 未关注  1 已关注 2 互相关注
func GetFollowingFlag(user_id, target_user_id string) (int, error) {

	params := map[string]string{"Selfuserid": user_id, "Targetuserid": target_user_id}
	var reply GetFlagRes

	err := UserRelationClient.Call("get_following_flag", params, &reply)
	if err != nil {
		Logger.Error(err.Error())
		return 0, err
	}

	return reply.Data, nil
}

// 批量查询
// 返回值 0 未关注  1 已关注 2 互相关注
func GetFollowingFlags(user_id string, target_user_ids []string) (map[string]int, error) {
	params := map[string]interface{}{"Selfuserid": user_id, "Targetuserids": target_user_ids}
	var reply GetFlagsRes
	err := UserRelationClient.Call("get_following_flags", params, &reply)
	if err != nil {
		Logger.Error(err.Error())
		return reply.Data, err
	}
	return reply.Data, err
}

func GetFollowingCount(user_id string) (int, error) {
	var user_id_data = map[string]string{
		"user_id": user_id,
	}
	var reply GetCountRes
	err := UserRelationClient.Call("get_following_count", user_id_data, &reply)
	if err != nil {
		Logger.Error(err.Error())
		return 0, err
	}

	return reply.Data, err
}

func GetFollowerCount(user_id string) (int, error) {
	var user_id_data = map[string]string{
		"user_id": user_id,
	}
	var reply GetCountRes
	err := UserRelationClient.Call("get_follower_count", user_id_data, &reply)
	if err != nil {
		Logger.Error(err.Error())
		return 0, err
	}

	return reply.Data, err
}

func GetFriendsCount(user_id string) (int, error) {
	var user_id_data = map[string]string{
		"user_id": user_id,
	}
	var reply GetCountRes
	err := UserRelationClient.Call("get_friends_count", user_id_data, &reply)
	if err != nil {
		Logger.Error(err.Error())
		return 0, err
	}

	return reply.Data, err
}
