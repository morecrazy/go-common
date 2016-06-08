package protocol

type GetAccountArgs struct {
	AccountId int64  `json:"account_id"`
	UserId    string `json:"user_id"`
}

type GetAccountReply struct {
	AccountId   int64  `json:"account_id"`
	UserId      string `json:"user_id"`
	Balance     int64  `json:"balance"`
	CashBalance int64  `json:"cash_balance"`
}

var CoinServerFuncMap map[string]string = map[string]string{
	"get_account": "CoinServerHandler.GetAccount",
}
