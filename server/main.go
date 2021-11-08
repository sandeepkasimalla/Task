package main

import (
	"log"
	"net/http"

	"github.com/go-chassis/openlog"
	"github.com/gorilla/mux"

	"Task/database"
	"Task/handlers"
	"Task/repository"
)

func GetHandler(dbname string) *handlers.Handler {
	repo := repository.UsersRepo{DbClient: database.GetClient(), DatabaseName: dbname}
	return &handlers.Handler{Repo: repo}
}
func main() {

	r := mux.NewRouter()
	err := database.Connect()
	if err != nil {
		openlog.Error(err.Error())
		return
	}
	h := GetHandler("Users")
	r.HandleFunc("/users/{id}", h.GetUser).Methods("GET")
	r.HandleFunc("/users", h.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")
	r.HandleFunc("/users", h.FetchAllUsers).Methods("GET")
	openlog.Info("Started listening at http://localhost:8070")
	log.Fatal(http.ListenAndServe(":8070", r))

}
