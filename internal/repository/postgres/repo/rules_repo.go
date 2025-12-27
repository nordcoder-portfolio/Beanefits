// internal/repository/postgres/repo/rules_repo.go
package repo

import (
	"context"
	"errors"
	"time"

	pg "Beanefits/internal/repository/postgres"
	pgdto "Beanefits/internal/repository/postgres/dto"
	"Beanefits/internal/repository/postgres/sqlc/gen"

	"github.com/jackc/pgx/v5"
)

type RulesRepo struct {
	q *gen.Queries
}

func NewRulesRepo(q *gen.Queries) *RulesRepo { return &RulesRepo{q: q} }

func (r *RulesRepo) GetEffectiveAt(ctx context.Context, db pg.DBTX, at time.Time) (pgdto.RulesetWithLevels, bool, error) {
	rs, err := r.q.GetRulesetEffectiveAt(ctx, db, timestamptz(at))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgdto.RulesetWithLevels{}, false, nil
		}
		return pgdto.RulesetWithLevels{}, false, err
	}

	levels, err := r.q.ListLevelRulesByRulesetID(ctx, db, rs.ID)
	if err != nil {
		return pgdto.RulesetWithLevels{}, false, err
	}

	return pgdto.RulesetWithLevels{
		Ruleset: mapRulesetEffective(rs),
		Levels:  mapLevelRules(levels),
	}, true, nil
}

func (r *RulesRepo) CreateRuleset(
	ctx context.Context,
	db pg.DBTX,
	effectiveFrom time.Time,
	baseRubPerPoint pgdto.Money,
	levels []pgdto.LevelRuleRow,
) (pgdto.RulesetWithLevels, error) {
	rs, err := r.q.InsertRuleset(ctx, db, gen.InsertRulesetParams{
		EffectiveFrom:   timestamptz(effectiveFrom),
		BaseRubPerPoint: baseRubPerPoint,
	})
	if err != nil {
		return pgdto.RulesetWithLevels{}, err
	}

	// вставляем уровни (IDs вернутся, но нам проще потом считать весь snapshot одним запросом)
	for _, lv := range levels {
		_, err := r.q.InsertLevelRule(ctx, db, gen.InsertLevelRuleParams{
			RulesetID:           rs.ID,
			LevelCode:           lv.LevelCode,
			ThresholdTotalSpend: lv.ThresholdTotalSpend,
			PercentEarn:         lv.PercentEarn,
		})
		if err != nil {
			return pgdto.RulesetWithLevels{}, err
		}
	}

	rs2, err := r.q.GetRulesetByID(ctx, db, rs.ID)
	if err != nil {
		return pgdto.RulesetWithLevels{}, err
	}
	levels2, err := r.q.ListLevelRulesByRulesetID(ctx, db, rs.ID)
	if err != nil {
		return pgdto.RulesetWithLevels{}, err
	}

	return pgdto.RulesetWithLevels{
		Ruleset: mapRulesetByID(rs2),
		Levels:  mapLevelRules(levels2),
	}, nil
}

func (r *RulesRepo) ListRulesets(ctx context.Context, db pg.DBTX, limit, offset int) ([]pgdto.RulesetWithLevels, error) {
	rows, err := r.q.ListRulesetsBase(ctx, db, gen.ListRulesetsBaseParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	out := make([]pgdto.RulesetWithLevels, 0, len(rows))
	for _, rs := range rows {
		levels, err := r.q.ListLevelRulesByRulesetID(ctx, db, rs.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, pgdto.RulesetWithLevels{
			Ruleset: mapRulesetBase(rs),
			Levels:  mapLevelRules(levels),
		})
	}
	return out, nil
}

// ---------- helpers ----------

func mapRulesetEffective(rw gen.GetRulesetEffectiveAtRow) pgdto.RulesetRow {
	return pgdto.RulesetRow{
		ID:              rw.ID,
		EffectiveFrom:   rw.EffectiveFrom.Time,
		BaseRubPerPoint: rw.BaseRubPerPoint,
		CreatedAt:       rw.CreatedAt.Time,
	}
}

func mapRulesetByID(rw gen.GetRulesetByIDRow) pgdto.RulesetRow {
	return pgdto.RulesetRow{
		ID:              rw.ID,
		EffectiveFrom:   rw.EffectiveFrom.Time,
		BaseRubPerPoint: rw.BaseRubPerPoint,
		CreatedAt:       rw.CreatedAt.Time,
	}
}

func mapRulesetBase(rw gen.ListRulesetsBaseRow) pgdto.RulesetRow {
	return pgdto.RulesetRow{
		ID:              rw.ID,
		EffectiveFrom:   rw.EffectiveFrom.Time,
		BaseRubPerPoint: rw.BaseRubPerPoint,
		CreatedAt:       rw.CreatedAt.Time,
	}
}

func mapLevelRules(rows []gen.LevelRule) []pgdto.LevelRuleRow {
	out := make([]pgdto.LevelRuleRow, 0, len(rows))
	for _, r := range rows {
		out = append(out, pgdto.LevelRuleRow{
			ID:                  r.ID,
			RulesetID:           r.RulesetID,
			LevelCode:           r.LevelCode,
			ThresholdTotalSpend: r.ThresholdTotalSpend,
			PercentEarn:         r.PercentEarn,
		})
	}
	return out
}
