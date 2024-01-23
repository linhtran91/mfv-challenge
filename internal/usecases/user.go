package usecases

type UserAccounts struct {
	ID       int64   `json:"id"`
	Username string  `json:"name"`
	Accounts []int64 `json:"account_ids"`
}
