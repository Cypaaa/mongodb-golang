package main

import (
	"atlas_front/app/router"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	srv := http.Server{
		Addr:         ":8080",
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 250 * time.Millisecond,
		Handler:      r,
	}

	router.Init(r)

	fmt.Println("Listening to", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
