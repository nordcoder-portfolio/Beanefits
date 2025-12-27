package client

import (
	"context"
	"fmt"
	"time"

	"Beanefits/internal/domain/errs"
	pgdto "Beanefits/internal/repository/postgres/dto"
	sdto "Beanefits/internal/service/dto"
	"Beanefits/internal/service/mapper"
)

func (s *Service) GetMe(ctx context.Context, userID int64) (sdto.ClientProfileOut, error) {
	start := time.Now()
	s.log.InfoContext(ctx, "client.get_me start", "userID", userID)

	u, err := s.requireActiveUser(ctx, userID)
	if err != nil {
		s.log.ErrorContext(ctx, "client.get_me failed", "ms", time.Since(start).Milliseconds(), "err", err)
		return sdto.ClientProfileOut{}, err
	}

	rs, err := s.requireRoles(ctx, userID)
	if err != nil {
		s.log.ErrorContext(ctx, "client.get_me failed", "ms", time.Since(start).Milliseconds(), "err", err)
		return sdto.ClientProfileOut{}, err
	}

	// sanity check: transport usually guarantees role, but keep guard
	if !hasRole(rs, pgdto.RoleClient) {
		e := errs.New(errs.CodeRolesNotFound, "client role not found")
		s.log.ErrorContext(ctx, "client.get_me failed", "ms", time.Since(start).Milliseconds(), "err", e)
		return sdto.ClientProfileOut{}, e
	}

	acc, err := s.requireAccount(ctx, userID)
	if err != nil {
		s.log.ErrorContext(ctx, "client.get_me failed", "ms", time.Since(start).Milliseconds(), "err", err)
		return sdto.ClientProfileOut{}, err
	}

	out := sdto.ClientProfileOut{
		User:    mapper.UserWithRoles(u, rs),
		Account: mapper.AccountBase(acc),
	}

	s.log.InfoContext(ctx, "client.get_me ok", "ms", time.Since(start).Milliseconds(), "userID", userID, "accountID", acc.ID, "level", fmt.Sprint(acc.LevelCode))
	return out, nil
}
