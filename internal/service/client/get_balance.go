package client

import (
	"context"
	"time"

	sdto "Beanefits/internal/service/dto"
	"Beanefits/internal/service/mapper"
)

func (s *Service) GetBalance(ctx context.Context, userID int64) (sdto.BalanceOut, error) {
	start := time.Now()
	s.log.InfoContext(ctx, "client.get_balance start", "userID", userID)

	if _, err := s.requireActiveUser(ctx, userID); err != nil {
		s.log.ErrorContext(ctx, "client.get_balance failed", "ms", time.Since(start).Milliseconds(), "err", err)
		return sdto.BalanceOut{}, err
	}

	acc, err := s.requireAccount(ctx, userID)
	if err != nil {
		s.log.ErrorContext(ctx, "client.get_balance failed", "ms", time.Since(start).Milliseconds(), "err", err)
		return sdto.BalanceOut{}, err
	}

	out := mapper.BalanceOut(acc, s.now())
	s.log.InfoContext(ctx, "client.get_balance ok", "ms", time.Since(start).Milliseconds(), "userID", userID, "accountID", acc.ID)
	return out, nil
}
