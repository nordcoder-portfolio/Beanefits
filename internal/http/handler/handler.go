package handler

import (
	"context"
	"net/http"

	"Beanefits/internal/api"
	"Beanefits/internal/service"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type Handler struct {
	api.Unimplemented

	authSvc    service.Auth
	clientSvc  service.Client
	cashierSvc service.Cashier
	adminSvc   service.Admin

	db      Pinger
	metrics http.Handler

	helpers Helpers
}

type Deps struct {
	Auth    service.Auth
	Client  service.Client
	Cashier service.Cashier
	Admin   service.Admin

	DB      Pinger
	Metrics http.Handler

	Helpers Helpers
}

func New(d Deps) *Handler {
	h := &Handler{
		authSvc:    d.Auth,
		clientSvc:  d.Client,
		cashierSvc: d.Cashier,
		adminSvc:   d.Admin,
		db:         d.DB,
		metrics:    d.Metrics,
		helpers:    d.Helpers,
	}

	if h.helpers.JSON == nil {
		h.helpers = DefaultHelpers()
	}

	return h
}
