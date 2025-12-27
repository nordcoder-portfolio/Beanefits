package handler

import (
	"net/http"

	"Beanefits/internal/api"
	sdto "Beanefits/internal/service/dto"
)

// GET /me
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.requireClient(w, r)
	if !ok {
		return
	}

	out, err := h.clientSvc.GetMe(r.Context(), userID)
	if err != nil {
		h.WriteServiceError(w, r, err)
		return
	}

	h.helpers.JSON(w, http.StatusOK, mapClientProfile(out))
}

// GET /me/balance
func (h *Handler) GetMeBalance(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.requireClient(w, r)
	if !ok {
		return
	}

	out, err := h.clientSvc.GetBalance(r.Context(), userID)
	if err != nil {
		h.WriteServiceError(w, r, err)
		return
	}

	h.helpers.JSON(w, http.StatusOK, mapBalance(out))
}

// GET /me/events
func (h *Handler) GetMeEvents(w http.ResponseWriter, r *http.Request, params api.GetMeEventsParams) {
	userID, ok := h.requireClient(w, r)
	if !ok {
		return
	}

	in := sdto.EventsIn{
		Limit:    derefLimit(params.Limit, 20),
		BeforeTs: params.BeforeTs,
	}

	out, err := h.clientSvc.GetEvents(r.Context(), userID, in)
	if err != nil {
		h.WriteServiceError(w, r, err)
		return
	}

	h.helpers.JSON(w, http.StatusOK, mapEventsPage(out))
}

func mapClientProfile(out sdto.ClientProfileOut) api.ClientProfile {
	return api.ClientProfile{
		User: api.User{
			Id:        out.User.ID,
			Phone:     api.Phone(out.User.Phone),
			Roles:     mapRolesToAPI(out.User.Roles),
			IsActive:  out.User.IsActive,
			CreatedAt: out.User.CreatedAt,
		},
		Account: api.Account{
			Id:              out.Account.ID,
			PublicCode:      api.PublicCode(out.Account.PublicCode),
			BalancePoints:   out.Account.BalancePoints,
			TotalSpendMoney: out.Account.TotalSpendMoney,
			LevelCode:       api.LevelCode(out.Account.LevelCode),
			CreatedAt:       out.Account.CreatedAt,
		},
	}
}
