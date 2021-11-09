package handlers

import (
	"Task/common"
	validator "Task/payloadvalidator"
	"Task/service"
	"encoding/json"
	"net/http"

	"github.com/go-chassis/openlog"
	"github.com/gorilla/mux"
)

type Handler struct {
	Service service.Service
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
	input := common.CreateUserInput{User: user}
	res := h.Service.CreateUser(input)
	w.WriteHeader(res.Status)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	openlog.Info("Got a request to fetch User")
	// set header.
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id := params["id"]
	input := common.FetchUserInput{ID: id}
	res := h.Service.GetUser(input)
	w.WriteHeader(res.Status)
	json.NewEncoder(w).Encode(res)
}

// FetchAllDatamodelsByPagenation function will helps to get the field by considering page number, size and filters.
func (h *Handler) FetchAllUsers(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	size := r.URL.Query().Get("size")
	filters := r.URL.Query().Get("filters")
	sort := r.URL.Query().Get("sort")
	input := common.FetchAllUsersInput{Page: page, Size: size, Filters: filters, Sort: sort}
	w.Header().Set("Content-Type", "application/json")
	res := h.Service.FetchAllUsers(input)
	w.WriteHeader(res.Status)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	openlog.Info("Got a request to Delete User")
	// set header.
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id := params["id"]
	input := common.DeleteUserInput{ID: id}
	res := h.Service.DeleteUser(input)
	w.WriteHeader(res.Status)
	json.NewEncoder(w).Encode(res)
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
	input := common.UpdateUserInput{ID: id, User: user}
	res := h.Service.UpdateUser(input)
	w.WriteHeader(res.Status)
	json.NewEncoder(w).Encode(res)
}
