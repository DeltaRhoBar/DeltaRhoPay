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
	updateResident http.Handler,
	allResidentsPage http.Handler,
	forceAddResident http.Handler, 
	addBeverage http.Handler, 
	removeBeverage http.Handler, 
	addOrder http.Handler, 
	getResidents http.Handler,
	getOrders http.Handler,
	debtPage http.Handler,
	checkout http.Handler,
	pay http.Handler,
) *Router {
	cr := chi.NewRouter()

	router := &Router {
		chiRouter: cr,
		authenticator: authenticator,
	}

	static := http.FileServer(http.Dir("web/static/"))

	router.setupRoutes(static, index, loginPage, admin, ordersPage, loginData, addResident, updateResident, allResidentsPage, forceAddResident, addBeverage, removeBeverage, addOrder, getResidents, getOrders, debtPage, checkout, pay)

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
	updateResident http.Handler,
	allResidentsPage http.Handler,
	forceAddResident http.Handler, 
	addBeverage http.Handler, 
	removeBeverage http.Handler, 
	addOrder http.Handler, 
	getResidents http.Handler,
	getOrders http.Handler,
	debtPage http.Handler,
	checkout http.Handler,
	pay http.Handler,
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
		cr.Handle("/updateResident", updateResident)
		cr.Handle("/residents", allResidentsPage)
		cr.Handle("/forceAddResident", forceAddResident)
		cr.Handle("/addBeverage", addBeverage)
		cr.Handle("/removeBeverage", removeBeverage)
		cr.Handle("/admin", admin)
		cr.Handle("/orders", ordersPage)
		cr.Handle("/getOrders", getOrders)
		cr.Handle("/debt", debtPage)
		cr.Handle("/checkout", checkout)
		cr.Handle("/pay", pay)
	})
}

func (r *Router) Handler() *chi.Mux {
	return r.chiRouter
}
