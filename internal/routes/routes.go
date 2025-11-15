package routes

import (
    "net/http"
    "github.com/go-chi/chi/v5"
)

type Router struct {
	chiRouter *chi.Mux
}

func NewRouter(index http.Handler, admin http.Handler, addResident http.Handler, forceAddResident http.Handler, addBeverage http.Handler, removeBeverage http.Handler, addDebt http.Handler) *Router {
	cr := chi.NewRouter()

	router := &Router {
		chiRouter: cr,
	}

	static := http.FileServer(http.Dir(".web"))

	router.setupMiddleware()
	router.setupRoutes(static, index, admin, addResident, forceAddResident, addBeverage, removeBeverage, addDebt)

	return router
}

func (r *Router) setupRoutes(static http.Handler, index http.Handler, admin http.Handler, addResident http.Handler, forceAddResident http.Handler, addBeverage http.Handler, removeBeverage http.Handler, addDebt http.Handler) {
	r.chiRouter.Handle("/static/*", http.StripPrefix("/static/", static))
	r.chiRouter.Handle("/", index)
	r.chiRouter.Handle("/admin", admin)
	r.chiRouter.Handle("/addResident", addResident)
	r.chiRouter.Handle("/forceAddResident", forceAddResident)
	r.chiRouter.Handle("/addBeverage", addBeverage)
	r.chiRouter.Handle("/removeBeverage", removeBeverage)
	r.chiRouter.Handle("/addDebt", addDebt)
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
