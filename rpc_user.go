package common

import (
	"backend/userprofile/userhobby"
)

type UserprofileGeoIndexArgs struct {
	Id string `json:"id"`
}

type UserprofileGeoIndexReply struct {
	Id        string              `json:"id"`
	HobbyList userhobby.HobbyList `json:"hobby_list"`
	Gender    string              `json:"gender"`
}

func GetGeoIndexData(upr UserprofileGeoIndexArgs) (resp UserprofileGeoIndexReply, err error) {

	reply := UserprofileGeoIndexReply{}
	err = UserProfileClient.Call("UserprofileHandler.GetGeoIndexData", &upr, &reply)
	if err != nil {
		Logger.Error("GetGeoIndexData rpc error", err)
	}
	return reply, err
}

