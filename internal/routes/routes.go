package routes

import (
    "github.com/go-chi/chi/v5"
)

type Router struct {
	chiRouter *chi.Mux
}

func NewRouter() *Router {
	cr := chi.NewRouter()

	router := &Router {
		chiRouter: cr,
	}
	return router
}

func (r *Router) setupRoutes() error {
	/*
	r.chiRouter.Route("/users", func(cr chi.Router) {
	cr.Get("/", function)
	}
	*/
	return nil
}

func (r *Router) Handler() *chi.Mux {
	return r.chiRouter
}
