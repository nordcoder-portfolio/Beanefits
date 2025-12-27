package dto

type AccountRow struct {
	ID              int64
	UserID          int64
	PublicCode      string
	BalancePoints   int
	TotalSpendMoney Money
	LevelCode       string
	CreatedAt       Ts
}
