package main

import (
	"log"
    "net/http"
	"deltapay/internal/handlers"
	"deltapay/internal/routes"
	"deltapay/internal/services"

)

func main() {
	database, err := services.NewSqlite();
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()


type PageRoom struct {
	Floor int
	Room int
}

type PageResident struct {
	Room PageRoom
	Name string
}
	indexHandler, err := handlers.NewIndexHandler(database)
	if err != nil {
		log.Fatal(err)
	}

	adminHandler, err := handlers.NewAdminHandler()
	if err != nil {
		log.Fatal(err)
	}

	addResident := handlers.NewAddResidentHandler(database)
	forceAddResident := handlers.NewForceAddResidentHandler(database)

	r := routes.NewRouter(indexHandler, adminHandler, addResident, forceAddResident)

    http.ListenAndServe(":8080", r.Handler())
}
