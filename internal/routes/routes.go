package routes

import (
	"deltapay/internal/middleware"
	"deltapay/internal/services"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	chiRouter *chi.Mux
	authenticator *services.Authenticator
}

func NewRouter(
	authenticator *services.Authenticator, 
	index http.Handler, 
	loginPage http.Handler, 
	admin http.Handler, 
	ordersPage http.Handler,
	loginData http.Handler, 
	addResident http.Handler, 
	forceAddResident http.Handler, 
	addBeverage http.Handler, 
	removeBeverage http.Handler, 
	addOrder http.Handler, 
	getResidents http.Handler,
) *Router {
	cr := chi.NewRouter()

	router := &Router {
		chiRouter: cr,
		authenticator: authenticator,
	}

	static := http.FileServer(http.Dir(".web"))

	router.setupRoutes(static, index, loginPage, admin, ordersPage, loginData, addResident, forceAddResident, addBeverage, removeBeverage, addOrder, getResidents)

	return router
}

func (r *Router) setupRoutes(
	static http.Handler, 
	index http.Handler, 
	loginPage http.Handler, 
	admin http.Handler, 
	ordersPage http.Handler,
	loginData http.Handler, 
	addResident http.Handler, 
	forceAddResident http.Handler, 
	addBeverage http.Handler, 
	removeBeverage http.Handler, 
	addOrder http.Handler, 
	getResidents http.Handler,
) {
	r.chiRouter.Handle("/static/*", http.StripPrefix("/static/", static))
	r.chiRouter.Handle("/", index)
	r.chiRouter.Handle("/login", loginPage)
	r.chiRouter.Handle("/loginData", loginData)
	r.chiRouter.Handle("/addOrder", addOrder)
	r.chiRouter.Handle("/getResidents", getResidents)
	r.chiRouter.Group(func(cr chi.Router) {
		cr.Use(middleware.Auth(r.authenticator))
		cr.Handle("/addResident", addResident)
		cr.Handle("/forceAddResident", forceAddResident)
		cr.Handle("/addBeverage", addBeverage)
		cr.Handle("/removeBeverage", removeBeverage)
		cr.Handle("/admin", admin)
		cr.Handle("/orders", ordersPage)
	})
}

func (r *Router) Handler() *chi.Mux {
	return r.chiRouter
}
