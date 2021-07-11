package main

import (
	"log"
	"net/http"

	"github.com/cenkayla/shorturl/models"
	"github.com/cenkayla/shorturl/service"
)

func main() {
	db, err := models.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	srv := service.InitService(db)
	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(":8080", srv)

}
