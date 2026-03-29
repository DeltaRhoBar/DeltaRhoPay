package main

import (
	"log"
    "net/http"
	"deltapay/internal/handlers"
	"deltapay/internal/routes"
	"deltapay/internal/services"
	"github.com/joho/godotenv"

)

func main() {
	godotenv.Load()

	database, err := services.NewSqlite();
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	authenticator := services.NewAuthenticator()

	indexHandler, err := handlers.NewIndexHandler(database)
	if err != nil {
		log.Fatal(err)
		return
	}

	loginPageHandler, err := handlers.NewLoginPageHandler()
	if err != nil {
		log.Fatal(err)
		return
	}

	adminHandler, err := handlers.NewAdminHandler(database)
	if err != nil {
		log.Fatal(err)
		return
	}

	residentsPageHandler, err := handlers.NewResidentsPageHandler(database)
	if err != nil {
		log.Fatal(err)
		return
	}

	ordersPageHandler, err := handlers.NewOrdersPageHandler(database)
	if err != nil {
		log.Fatal(err)
		return
	}

	debtPageHandler, err := handlers.NewDebtPageHandler(database)
	if err != nil {
		log.Fatal(err)
		return
	}


	loginData := handlers.NewLoginHandler(authenticator)
	addResident := handlers.NewAddResidentHandler(database)
	updateResident := handlers.NewUpdateResidentHandler(database)
	forceAddResident := handlers.NewForceAddResidentHandler(database)
	addBeverage := handlers.NewAddBeverageHandler(database)
	removeBeverage := handlers.NewRemoveBeverageHandler(database)
	addOrder:= handlers.NewAddOrderHandler(database)
	getResidents := handlers.NewGetResidentHandler(database)
	getOrders := handlers.NewGetOrdersHandler(database)
	checkout := handlers.NewCheckoutHandler(database)
	pay := handlers.NewPayHandler(database)


	r := routes.NewRouter(
		authenticator, 
		indexHandler, 
		loginPageHandler, 
		adminHandler, 
		ordersPageHandler, 
		loginData, 
		addResident, updateResident, residentsPageHandler, forceAddResident, addBeverage, removeBeverage, addOrder, getResidents, getOrders, debtPageHandler, checkout,
		pay)

    http.ListenAndServe(":8080", r.Handler())
}
