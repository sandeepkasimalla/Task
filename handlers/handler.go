package handlers

import (
	validator "Task/payloadvalidator"
	"Task/repository"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chassis/openlog"
	"github.com/gorilla/mux"
)

type Handler struct {
	Repo repository.UsersRepo
}

type Response struct {
	Msg    string      `json:"_msg"`
	Status int         `json:"_status"`
	Data   interface{} `json:"data"`
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	openlog.Info("Got a request to create User")
	w.Header().Set("Content-Type", "application/json")

	user := make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&user)
	valres, err := validator.ValidatePaylaod("./../payloadschemas/user.json", user)
	if err != nil {
		openlog.Error(err.Error())
		response := Response{Msg: err.Error(), Data: valres, Status: 400}
		json.NewEncoder(w).Encode(response)
		return
	}
	email := user["email"].(string)
	code, err := h.Repo.IsEmailExists(email)
	fmt.Println(err, " ++++++++++++++++++++++++++++++++++")
	if err != nil {
		openlog.Error(err.Error())
		response := Response{Msg: err.Error(), Data: nil, Status: code}
		json.NewEncoder(w).Encode(response)
		return
	}
	res, code, err := h.Repo.Insert(user)
	if err != nil {
		response := Response{Msg: err.Error(), Data: nil, Status: code}
		json.NewEncoder(w).Encode(response)
		return
	}
	response := Response{Msg: "User inserted successfully", Data: res, Status: 201}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	openlog.Info("Got a request to fetch User")
	// set header.
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id := params["id"]
	res, code, err := h.Repo.Find(id)
	if err != nil {
		openlog.Error(err.Error())
		response := Response{Msg: err.Error(), Data: nil, Status: code}
		json.NewEncoder(w).Encode(response)
		return
	}
	response := Response{Msg: "User Fetched successfully", Data: res, Status: 200}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) FetchAllUsers(w http.ResponseWriter, r *http.Request) {
	openlog.Info("Got a request to fetch User")
	// set header.
	w.Header().Set("Content-Type", "application/json")
	filters := r.URL.Query().Get("filters")

	var filter = make(map[string]interface{})
	if filters != "" {
		bytes := []byte(filters)
		json.Unmarshal(bytes, &filter)
	}
	res, code, err := h.Repo.FindByFilters(filter, nil)
	if err != nil {
		openlog.Error(err.Error())
		response := Response{Msg: err.Error(), Data: nil, Status: code}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{Msg: "User Fetched successfully", Data: res, Status: 200}
	json.NewEncoder(w).Encode(response)
}
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	openlog.Info("Got a request to Delete User")
	// set header.
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id := params["id"]
	res, code, err := h.Repo.Delete(id)
	if err != nil {
		openlog.Error(err.Error())
		response := Response{Msg: err.Error(), Data: nil, Status: code}
		json.NewEncoder(w).Encode(response)
		return
	}
	response := Response{Msg: "User deleted successfully", Data: res, Status: 200}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	openlog.Info("Got a request to Update User")
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id := params["id"]

	user := make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&user)
	valres, err := validator.ValidatePaylaod("./../payloadschemas/updateuser.json", user)
	if err != nil {
		openlog.Error(err.Error())
		response := Response{Msg: err.Error(), Data: valres, Status: 400}
		json.NewEncoder(w).Encode(response)
		return
	}
	res, code, err := h.Repo.Find(id)
	if err != nil {
		openlog.Error(err.Error())
		response := Response{Msg: err.Error(), Data: nil, Status: code}
		json.NewEncoder(w).Encode(response)
		return
	}
	email, ok := user["email"].(string)
	if ok {
		code, err = h.Repo.IsEmailExists(email)
		if err != nil {
			openlog.Error(err.Error())
			response := Response{Msg: err.Error(), Data: nil, Status: code}
			json.NewEncoder(w).Encode(response)
			return
		}
	}
	res, code, err = h.Repo.FindAndUpdate(id, user)
	if err != nil {
		openlog.Error(err.Error())
		response := Response{Msg: err.Error(), Data: nil, Status: code}
		json.NewEncoder(w).Encode(response)
		return
	}
	response := Response{Msg: "User Updated successfully", Data: res, Status: 201}
	json.NewEncoder(w).Encode(response)
}
