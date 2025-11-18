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
	}

	loginPageHandler, err := handlers.NewLoginPageHandler()
	if err != nil {
		log.Fatal(err)
	}

	adminHandler, err := handlers.NewAdminHandler(database)
	if err != nil {
		log.Fatal(err)
	}

	ordersPageHandler, err := handlers.NewOrdersPageHandler(database)
	if err != nil {
		log.Fatal(err)
	}

	loginData := handlers.NewLoginHandler(authenticator)
	addResident := handlers.NewAddResidentHandler(database)
	forceAddResident := handlers.NewForceAddResidentHandler(database)
	addBeverage := handlers.NewAddBeverageHandler(database)
	removeBeverage := handlers.NewRemoveBeverageHandler(database)
	addOrder:= handlers.NewAddOrderHandler(database)
	getResidents := handlers.NewGetResidentHandler(database)

	r := routes.NewRouter(authenticator, indexHandler, loginPageHandler, adminHandler, ordersPageHandler, loginData, addResident, forceAddResident, addBeverage, removeBeverage, addOrder, getResidents)

    http.ListenAndServe(":8080", r.Handler())
}
