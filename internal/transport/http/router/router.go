package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	User UserRouter
}

type Handler struct {
	User UserHandler
}

func NewRouter(h Handler) http.Handler {
	r := chi.NewRouter()

	router := &Router{
		User: *NewUserRouter(h.User),
	}

	router.User.RegisterRoutes(r)

	return r
}
