package mapper

import (
	"time"

	pgdto "Beanefits/internal/repository/postgres/dto"
	sdto "Beanefits/internal/service/dto"
)

// Roles maps repository role codes to service role codes.
// Unknown roles are passed through as-is (best-effort).
func Roles(in []pgdto.RoleCode) []sdto.RoleCode {
	out := make([]sdto.RoleCode, 0, len(in))
	for _, r := range in {
		switch r {
		case pgdto.RoleClient:
			out = append(out, sdto.RoleClient)
		case pgdto.RoleCashier:
			out = append(out, sdto.RoleCashier)
		case pgdto.RoleAdmin:
			out = append(out, sdto.RoleAdmin)
		default:
			out = append(out, sdto.RoleCode(r))
		}
	}
	return out
}

func UserWithRoles(u pgdto.UserRow, roles []pgdto.RoleCode) sdto.UserWithRoles {
	return sdto.UserWithRoles{
		UserBase: sdto.UserBase{
			ID:        u.ID,
			Phone:     u.Phone,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt,
		},
		Roles: Roles(roles),
	}
}

// UserWithRolesFrom is a convenience mapper for pgdto.UserWithRoles.
func UserWithRolesFrom(u pgdto.UserWithRoles) sdto.UserWithRoles {
	return UserWithRoles(u.UserRow, u.Roles)
}

func AccountBase(a pgdto.AccountRow) sdto.AccountBase {
	return sdto.AccountBase{
		ID:              a.ID,
		PublicCode:      a.PublicCode,
		BalancePoints:   a.BalancePoints,
		TotalSpendMoney: MoneyFixed2(a.TotalSpendMoney),
		LevelCode:       a.LevelCode,
		CreatedAt:       a.CreatedAt,
	}
}

func BalanceOut(a pgdto.AccountRow, asOf time.Time) sdto.BalanceOut {
	return sdto.BalanceOut{
		AccountID:       a.ID,
		BalancePoints:   a.BalancePoints,
		TotalSpendMoney: MoneyFixed2(a.TotalSpendMoney),
		LevelCode:       a.LevelCode,
		AsOf:            asOf,
	}
}

func EventOut(e pgdto.EventRow) sdto.EventOut {
	var amount *string
	if e.AmountMoney != nil {
		s := MoneyPtrFixed2(e.AmountMoney)
		amount = &s
	}

	var typ sdto.EventType
	switch e.Type {
	case pgdto.EventEarn:
		typ = sdto.EventEarn
	case pgdto.EventSpend:
		typ = sdto.EventSpend
	default:
		typ = sdto.EventType(e.Type)
	}

	return sdto.EventOut{
		ID:           e.ID,
		AccountID:    e.AccountID,
		Type:         typ,
		DeltaPoints:  e.DeltaPoints,
		BalanceAfter: e.BalanceAfter,
		AmountMoney:  amount,
		RulesetID:    e.RulesetID,
		ActorUserID:  e.ActorUserID,
		Ts:           e.Ts,
	}
}

func RulesetOut(r pgdto.RulesetWithLevels) sdto.RulesetOut {
	levels := make([]sdto.LevelRuleOut, 0, len(r.Levels))
	for _, lv := range r.Levels {
		levels = append(levels, sdto.LevelRuleOut{
			ID:                  lv.ID,
			LevelCode:           lv.LevelCode,
			ThresholdTotalSpend: MoneyFixed2(lv.ThresholdTotalSpend),
			PercentEarn:         MoneyFixed2(lv.PercentEarn),
		})
	}

	return sdto.RulesetOut{
		ID:              r.Ruleset.ID,
		EffectiveFrom:   r.Ruleset.EffectiveFrom,
		BaseRubPerPoint: MoneyFixed2(r.Ruleset.BaseRubPerPoint),
		Levels:          levels,
		CreatedAt:       r.Ruleset.CreatedAt,
	}
}

// AuthOut is a small convenience mapper for auth usecases.
func AuthOut(token string, u pgdto.UserWithRoles, acc pgdto.AccountRow) sdto.AuthOut {
	return sdto.AuthOut{
		AccessToken: token,
		User:        UserWithRolesFrom(u),
		Account:     AccountBase(acc),
	}
}

// MoneyFixed2 formats repository money (decimal) as fixed-2 string for API.
func MoneyFixed2(m pgdto.Money) string {
	return m.StringFixed(2)
}

func MoneyPtrFixed2(m *pgdto.Money) string {
	if m == nil {
		return "0.00"
	}
	return m.StringFixed(2)
}
