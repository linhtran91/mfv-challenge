package usecases

type Account struct {
	ID      int64   `json:"id"`
	UserID  int64   `json:"user_id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}
