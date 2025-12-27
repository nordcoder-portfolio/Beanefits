package handler

import (
	"net/http"
	"time"

	"Beanefits/internal/api"
	"Beanefits/internal/http/middleware"
	"Beanefits/internal/service/dto"
)

// ===== Admin: Users =====

func (h *Handler) GetAdminUsers(w http.ResponseWriter, r *http.Request, params api.GetAdminUsersParams) {
	_, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}

	in := dto.ListUsersIn{
		Limit:  derefLimit(params.Limit, 20),
		Offset: derefInt(params.Offset, 0),
		Q:      derefString(params.Q, ""),
	}

	out, err := h.adminSvc.ListUsers(r.Context(), in)
	if err != nil {
		h.WriteServiceError(w, r, err)
		return
	}

	resp := api.UsersPage{
		Items: make([]api.User, 0, len(out.Items)),
		Total: out.Total,
	}

	for _, u := range out.Items {
		resp.Items = append(resp.Items, api.User{
			Id:        u.ID,
			Phone:     u.Phone,
			Roles:     mapRolesToAPI(u.Roles),
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt,
		})
	}

	h.helpers.JSON(w, http.StatusOK, resp)
}

func (h *Handler) DeleteAdminUsersUserId(w http.ResponseWriter, r *http.Request, userId int64) {
	_, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}

	if err := h.adminSvc.DeactivateUser(r.Context(), userId); err != nil {
		h.WriteServiceError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ===== Admin: Rulesets =====

func (h *Handler) GetAdminRulesets(w http.ResponseWriter, r *http.Request, params api.GetAdminRulesetsParams) {
	_, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}

	in := dto.ListRulesetsIn{
		Limit:  derefLimit(params.Limit, 20),
		Offset: derefInt(params.Offset, 0),
	}

	out, err := h.adminSvc.ListRulesets(r.Context(), in)
	if err != nil {
		h.WriteServiceError(w, r, err)
		return
	}

	resp := api.RulesetsPage{
		Items: make([]api.Ruleset, 0, len(out.Items)),
		Total: out.Total,
	}

	for _, rs := range out.Items {
		resp.Items = append(resp.Items, mapRulesetToAPI(rs))
	}

	h.helpers.JSON(w, http.StatusOK, resp)
}

func (h *Handler) PostAdminRulesets(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}

	var req api.PostAdminRulesetsJSONRequestBody
	if err := DecodeJSON(r, &req); err != nil {
		detail := err.Error()
		WriteProblem(
			w,
			http.StatusUnprocessableEntity,
			"Validation error",
			&detail,
			ptr("INVALID_RULESET"),
			instanceFromRequest(r),
		)
		return
	}

	in := dto.CreateRulesetIn{
		EffectiveFrom:   req.EffectiveFrom,
		BaseRubPerPoint: req.BaseRubPerPoint,
		Levels:          make([]dto.LevelRuleIn, 0, len(req.Levels)),
	}
	for _, lvl := range req.Levels {
		in.Levels = append(in.Levels, dto.LevelRuleIn{
			LevelCode:           string(lvl.LevelCode),
			ThresholdTotalSpend: lvl.ThresholdTotalSpend,
			PercentEarn:         lvl.PercentEarn,
		})
	}

	out, err := h.adminSvc.CreateRuleset(r.Context(), actorUserID, in)
	if err != nil {
		h.WriteServiceError(w, r, err)
		return
	}

	h.helpers.JSON(w, http.StatusCreated, mapRulesetToAPI(out))
}

func (h *Handler) GetAdminRulesetsCurrent(w http.ResponseWriter, r *http.Request) {
	_, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}

	out, err := h.adminSvc.GetCurrentRuleset(r.Context(), time.Now().UTC())
	if err != nil {
		h.WriteServiceError(w, r, err)
		return
	}

	h.helpers.JSON(w, http.StatusOK, mapRulesetToAPI(out))
}

// ===== RBAC =====

func (h *Handler) requireAdmin(w http.ResponseWriter, r *http.Request) (int64, bool) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		detail := "missing auth context"
		WriteProblem(w, http.StatusUnauthorized, "Unauthorized", &detail, ptr("UNAUTHORIZED"), instanceFromRequest(r))
		return 0, false
	}

	if !hasRole(claims.Roles, string(api.ADMIN)) {
		detail := "admin role required"
		WriteProblem(w, http.StatusForbidden, "Forbidden", &detail, ptr("FORBIDDEN"), instanceFromRequest(r))
		return 0, false
	}

	return claims.UserID, true
}

func hasRole(roles []string, want string) bool {
	for _, r := range roles {
		if r == want {
			return true
		}
	}
	return false
}

// ===== Mapping =====

func mapRolesToAPI(in []dto.RoleCode) []api.RoleCode {
	out := make([]api.RoleCode, 0, len(in))
	for _, r := range in {
		out = append(out, api.RoleCode(r))
	}
	return out
}

func mapRulesetToAPI(in dto.RulesetOut) api.Ruleset {
	levels := make([]api.LevelRule, 0, len(in.Levels))
	for _, lvl := range in.Levels {
		levels = append(levels, api.LevelRule{
			Id:                  lvl.ID,
			LevelCode:           lvl.LevelCode,
			ThresholdTotalSpend: lvl.ThresholdTotalSpend,
			PercentEarn:         lvl.PercentEarn,
		})
	}
	return api.Ruleset{
		Id:              in.ID,
		EffectiveFrom:   in.EffectiveFrom,
		BaseRubPerPoint: in.BaseRubPerPoint,
		Levels:          levels,
		CreatedAt:       in.CreatedAt,
	}
}

// ===== Param helpers =====

func derefLimit(p *api.LimitParam, def int) int {
	if p == nil {
		return def
	}
	return int(*p)
}

func derefInt(p *int, def int) int {
	if p == nil {
		return def
	}
	return *p
}

func derefString(p *string, def string) string {
	if p == nil {
		return def
	}
	return *p
}
