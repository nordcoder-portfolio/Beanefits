package handler

import (
	"github.com/google/uuid"
	"net/http"

	"Beanefits/internal/api"
	"Beanefits/internal/http/middleware"
	sdto "Beanefits/internal/service/dto"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// GET /cashier/accounts/by-code/{publicCode}
func (h *Handler) GetCashierAccountsByCodePublicCode(w http.ResponseWriter, r *http.Request, publicCode api.PublicCode) {
	_, ok := h.requireCashier(w, r)
	if !ok {
		return
	}

	out, err := h.cashierSvc.LookupAccountByCode(r.Context(), string(publicCode))
	if err != nil {
		h.WriteServiceError(w, r, err)
		return
	}

	h.helpers.JSON(w, http.StatusOK, mapCashierAccountSummary(out))
}

// GET /cashier/accounts/by-code/{publicCode}/events
func (h *Handler) GetCashierAccountsByCodePublicCodeEvents(
	w http.ResponseWriter,
	r *http.Request,
	publicCode api.PublicCode,
	params api.GetCashierAccountsByCodePublicCodeEventsParams,
) {
	_, ok := h.requireCashier(w, r)
	if !ok {
		return
	}

	in := sdto.EventsIn{
		Limit:    derefLimit(params.Limit, 20),
		BeforeTs: params.BeforeTs,
	}

	out, err := h.cashierSvc.GetAccountEventsByCode(r.Context(), string(publicCode), in)
	if err != nil {
		h.WriteServiceError(w, r, err)
		return
	}

	h.helpers.JSON(w, http.StatusOK, mapEventsPage(out))
}

// POST /cashier/earn
func (h *Handler) PostCashierEarn(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := h.requireCashier(w, r)
	if !ok {
		return
	}

	var req api.PostCashierEarnJSONRequestBody
	if err := DecodeJSON(r, &req); err != nil {
		detail := err.Error()
		WriteProblem(w, http.StatusUnprocessableEntity, "Validation error", &detail, ptr("INVALID_REQUEST"), instanceFromRequest(r))
		return
	}

	out, err := h.cashierSvc.Earn(r.Context(), actorUserID, sdto.EarnIn{
		OperationID: req.OperationId.String(),
		PublicCode:  string(req.PublicCode),
		AmountMoney: req.AmountMoney,
		Ts:          req.Ts,
	})
	if err != nil {
		h.WriteServiceError(w, r, err)
		return
	}

	h.helpers.JSON(w, http.StatusOK, mapOperationResult(out))
}

// POST /cashier/spend
func (h *Handler) PostCashierSpend(w http.ResponseWriter, r *http.Request) {
	actorUserID, ok := h.requireCashier(w, r)
	if !ok {
		return
	}

	var req api.PostCashierSpendJSONRequestBody
	if err := DecodeJSON(r, &req); err != nil {
		detail := err.Error()
		WriteProblem(w, http.StatusUnprocessableEntity, "Validation error", &detail, ptr("INVALID_REQUEST"), instanceFromRequest(r))
		return
	}

	out, err := h.cashierSvc.Spend(r.Context(), actorUserID, sdto.SpendIn{
		OperationID:  req.OperationId.String(),
		PublicCode:   string(req.PublicCode),
		AmountPoints: req.AmountPoints,
		Ts:           req.Ts,
	})
	if err != nil {
		h.WriteServiceError(w, r, err)
		return
	}

	h.helpers.JSON(w, http.StatusOK, mapOperationResult(out))
}

// ===== RBAC =====

func (h *Handler) requireCashier(w http.ResponseWriter, r *http.Request) (int64, bool) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		detail := "missing auth context"
		WriteProblem(w, http.StatusUnauthorized, "Unauthorized", &detail, ptr("UNAUTHORIZED"), instanceFromRequest(r))
		return 0, false
	}

	// строго CASHIER (как в OpenAPI description). Если хочешь, легко расширим до CASHIER||ADMIN.
	if !hasRole(claims.Roles, string(api.CASHIER)) {
		detail := "cashier role required"
		WriteProblem(w, http.StatusForbidden, "Forbidden", &detail, ptr("FORBIDDEN"), instanceFromRequest(r))
		return 0, false
	}

	return claims.UserID, true
}

// ===== Mapping =====

func mapCashierAccountSummary(out sdto.AccountOut) api.CashierAccountSummary {
	return api.CashierAccountSummary{
		AccountId:       out.ID,
		PublicCode:      api.PublicCode(out.PublicCode),
		BalancePoints:   out.BalancePoints,
		TotalSpendMoney: out.TotalSpendMoney,
		LevelCode:       api.LevelCode(out.LevelCode),
	}
}

func mapEventsPage(out sdto.EventsOut) api.EventsPage {
	items := make([]api.Event, 0, len(out.Items))
	for _, e := range out.Items {
		items = append(items, mapEvent(e))
	}

	return api.EventsPage{
		Items:        items,
		NextBeforeTs: out.NextBeforeTs,
	}
}

func mapEvent(e sdto.EventOut) api.Event {
	return api.Event{
		Id:           e.ID,
		AccountId:    e.AccountID,
		Type:         api.EventType(e.Type),
		DeltaPoints:  e.DeltaPoints,
		BalanceAfter: e.BalanceAfter,
		AmountMoney:  e.AmountMoney,
		RulesetId:    e.RulesetID,
		ActorUserId:  e.ActorUserID,
		Ts:           e.Ts,
	}
}

func mapOperationResult(out sdto.OperationOut) api.OperationResult {
	replay := out.IdempotentReplay
	uid, _ := uuid.Parse(out.OperationID)
	return api.OperationResult{
		OperationId:      openapi_types.UUID(uid),
		OpType:           api.OperationType(out.OpType),
		Event:            mapEvent(out.Event),
		Balance:          mapBalance(out.Balance),
		IdempotentReplay: &replay,
	}
}

func mapBalance(b sdto.BalanceOut) api.BalanceResponse {
	return api.BalanceResponse{
		AccountId:       b.AccountID,
		BalancePoints:   b.BalancePoints,
		TotalSpendMoney: b.TotalSpendMoney,
		LevelCode:       api.LevelCode(b.LevelCode),
		AsOf:            b.AsOf,
	}
}
