package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserRouter struct {
	handler UserHandler
}

type UserHandler interface {
	Registration(rw http.ResponseWriter, r *http.Request)
	Login(rw http.ResponseWriter, r *http.Request)
}

func NewUserRouter(h UserHandler) *UserRouter {
	return &UserRouter{handler: h}
}

func (u *UserRouter) RegisterRoutes(r chi.Router) {
	r.Route("/api/user", func(r chi.Router) {
		r.Post(("/register"), u.handler.Registration)
		r.Post("/login", u.handler.Login)
	})
}
