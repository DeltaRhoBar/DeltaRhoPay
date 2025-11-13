package main

import (
    "net/http"
	"deltapay/internal/routes"
)

func main() {
	r := routes.NewRouter()
    http.ListenAndServe(":8080", r.Handler())
}
