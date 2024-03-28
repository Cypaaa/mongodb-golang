package handler

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type pollData struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func CreatePollGET(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("public/template/poll/index.gotmpl")
	tmpl.Execute(w, pollData{})
}

func CreatePollPOST(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	b, _ := io.ReadAll(r.Body)

	// parse req, check if nil or name is empty
	err := json.Unmarshal(b, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Name) < 1 {
		http.Error(w, "empty name provided", http.StatusBadRequest)
		return
	}

	resp, err := http.Post("http://localhost:3000/poll", "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Println("[ERROR] (CreatePoll)", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	b, _ = io.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	w.Write(b)
}

func GetPollGET(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	resp, err := http.Get("http://localhost:3000/poll/" + id)
	if err != nil {
		log.Println("[ERROR] (CreatePoll)", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != 200 {
		http.Redirect(w, r, "http://localhost:8080/poll/new", http.StatusTemporaryRedirect)
	}

	var data pollData
	json.NewDecoder(resp.Body).Decode(&data)

	tmpl, _ := template.ParseFiles("public/template/poll/index.gotmpl")
	tmpl.Execute(w, data)
}

func UpdatePollPUT(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id := mux.Vars(r)["id"]
	req := pollData{
		ID: id,
	}

	b, _ := io.ReadAll(r.Body)

	// parse req, check if nil or name is empty
	err := json.Unmarshal(b, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Name) < 1 {
		http.Error(w, "empty name provided", http.StatusBadRequest)
		return
	}

	jsonData, _ := json.Marshal(req)
	request, _ := http.NewRequest(http.MethodPut, "http://localhost:3000/poll", bytes.NewBuffer(jsonData))
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println("[ERROR] (UpdatePoll)", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	b, _ = io.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	w.Write(b)
}

func DeletePollDELETE(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id := mux.Vars(r)["id"]
	req := pollData{
		ID: id,
	}

	b, _ := io.ReadAll(r.Body)

	// parse req, check if nil or name is empty
	err := json.Unmarshal(b, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, _ := json.Marshal(req)
	request, _ := http.NewRequest(http.MethodDelete, "http://localhost:3000/poll", bytes.NewBuffer(jsonData))
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println("[ERROR] (DeletePoll)", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	b, _ = io.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	w.Write(b)
}
