package handler

import (
	"net/http"

	"Beanefits/internal/api"
	sdto "Beanefits/internal/service/dto"
)

// POST /auth/register
func (h *Handler) PostAuthRegister(w http.ResponseWriter, r *http.Request) {
	var req api.PostAuthRegisterJSONRequestBody
	if err := DecodeJSON(r, &req); err != nil {
		detail := err.Error()
		WriteProblem(
			w,
			http.StatusUnprocessableEntity,
			"Validation error",
			&detail,
			ptr("INVALID_REQUEST"),
			instanceFromRequest(r),
		)
		return
	}

	out, err := h.authSvc.RegisterClient(r.Context(), sdto.RegisterIn{
		Phone:    string(req.Phone),
		Password: req.Password,
	})
	if err != nil {
		h.WriteServiceError(w, r, err)
		return
	}

	h.helpers.JSON(w, http.StatusCreated, mapAuthOutToAPI(out))
}

// POST /auth/login
func (h *Handler) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	var req api.PostAuthLoginJSONRequestBody
	if err := DecodeJSON(r, &req); err != nil {
		detail := err.Error()
		WriteProblem(
			w,
			http.StatusUnprocessableEntity,
			"Validation error",
			&detail,
			ptr("INVALID_REQUEST"),
			instanceFromRequest(r),
		)
		return
	}

	out, err := h.authSvc.Login(r.Context(), sdto.LoginIn{
		Phone:    string(req.Phone),
		Password: req.Password,
	})
	if err != nil {
		h.WriteServiceError(w, r, err)
		return
	}

	h.helpers.JSON(w, http.StatusOK, mapAuthOutToAPI(out))
}

func mapAuthOutToAPI(out sdto.AuthOut) api.AuthResponse {
	user := api.User{
		Id:        out.User.ID,
		Phone:     api.Phone(out.User.Phone),
		Roles:     mapRolesToAPI(out.User.Roles),
		IsActive:  out.User.IsActive,
		CreatedAt: out.User.CreatedAt,
	}

	acc := api.Account{
		Id:              out.Account.ID,
		PublicCode:      api.PublicCode(out.Account.PublicCode),
		BalancePoints:   out.Account.BalancePoints,
		TotalSpendMoney: out.Account.TotalSpendMoney,
		LevelCode:       api.LevelCode(out.Account.LevelCode),
		CreatedAt:       out.Account.CreatedAt,
	}

	return api.AuthResponse{
		AccessToken: out.AccessToken,
		User:        user,
		Account:     &acc,
	}
}
