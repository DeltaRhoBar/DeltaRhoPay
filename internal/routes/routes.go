package routes

import (
    "net/http"
    "github.com/go-chi/chi/v5"
)

type Router struct {
	chiRouter *chi.Mux
}

func NewRouter(index http.Handler, admin http.Handler, addResident http.Handler) *Router {
	cr := chi.NewRouter()

	router := &Router {
		chiRouter: cr,
	}

	static := http.FileServer(http.Dir(".web"))

	router.setupMiddleware()
	router.setupRoutes(static, index, admin, addResident)

	return router
}

func (r *Router) setupRoutes(static http.Handler, index http.Handler, admin http.Handler, addResident http.Handler) error {
	r.chiRouter.Handle("/static/*", http.StripPrefix("/static/", static))
	r.chiRouter.Handle("/", index)
	r.chiRouter.Handle("/admin", admin)
	r.chiRouter.Handle("/addResident", addResident)
	return nil
}

func (r *Router) setupMiddleware() error {
	/*
	r.chiRouter.Use(middleware)
	*/
	return nil
}

func (r *Router) Handler() *chi.Mux {
	return r.chiRouter
}
