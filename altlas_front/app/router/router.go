package router

import (
	"atlas_front/app/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func Init(r *mux.Router) {
	r.Handle("/", http.RedirectHandler("/poll/new", http.StatusPermanentRedirect))
	r.HandleFunc("/poll/new", handler.CreatePollGET).Methods("GET")
	r.HandleFunc("/poll/new", handler.CreatePollPOST).Methods("POST")
	r.HandleFunc("/poll/{id}", handler.GetPollGET).Methods("GET")
	r.HandleFunc("/poll/{id}", handler.UpdatePollPUT).Methods("PUT")
	r.HandleFunc("/poll/{id}", handler.DeletePollDELETE).Methods("DELETE")
}
