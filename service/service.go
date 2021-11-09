package service

import (
	"Task/common"
	"Task/repository"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-chassis/openlog"
)

type Service struct {
	Repo repository.UsersRepo
}

func (h *Service) CreateUser(input common.CreateUserInput) common.Response {
	email := input.User["email"].(string)
	code, err := h.Repo.IsEmailExists(email)
	fmt.Println(err, " ++++++++++++++++++++++++++++++++++")
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}
	}
	res, code, err := h.Repo.Insert(input.User)
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}
	}
	return common.Response{Msg: "User inserted successfully", Data: res, Status: 201}
}

func (h *Service) GetUser(input common.FetchUserInput) common.Response {

	res, code, err := h.Repo.Find(input.ID)
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}
	}
	return common.Response{Msg: "User Fetched successfully", Data: res, Status: 200}
}

//CheckAndInitializePageDetails will initilize page,and it's size
func CheckAndInitializePageDetails(page, size string) (int64, int64) {
	pageNo, err := strconv.ParseInt(page, 10, 64)
	if err != nil || pageNo < 1 {
		pageNo = -1
	}

	var pageSize int64
	pageSize, err = strconv.ParseInt(size, 10, 64)
	if err != nil || pageSize < 1 {
		pageSize = 5
	}
	return pageNo, pageSize
}

// FetchAllDatamodelsByPagenation function will helps to get the field by considering page number, size and filters.
func (h *Service) FetchAllUsers(input common.FetchAllUsersInput) common.Response {
	var filter = make(map[string]interface{})
	if input.Filters != "" {
		bytes := []byte(input.Filters)
		json.Unmarshal(bytes, &filter)
	}
	var sortorder = make(map[string]interface{})
	if input.Sort != "" {
		bytes := []byte(input.Sort)
		_ = json.Unmarshal(bytes, &sortorder)
	}
	var res = make([]map[string]interface{}, 0)
	var err error
	var code int
	pageNo, pageSize := CheckAndInitializePageDetails(input.Page, input.Size)
	if pageNo < 0 {
		res, code, err = h.Repo.FindByFilters(filter, sortorder)
		if err != nil {
			openlog.Error(err.Error())
			return common.Response{Msg: err.Error(), Data: nil, Status: code}
		}
	} else {
		res, code, err = h.Repo.FindByFiltersAndPagenation(pageNo, pageSize, filter, sortorder)
		if err != nil {
			openlog.Error(err.Error())
			return common.Response{Msg: err.Error(), Data: nil, Status: code}

		}
	}
	return common.Response{Msg: "All Users Fetched successfully", Data: res, Status: 200}
}

func (h *Service) DeleteUser(input common.DeleteUserInput) common.Response {
	openlog.Info("Got a request to Delete User")
	// set header.
	res, code, err := h.Repo.Delete(input.ID)
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}

	}
	return common.Response{Msg: "User deleted successfully", Data: res, Status: 200}
}

func (h *Service) UpdateUser(input common.UpdateUserInput) common.Response {
	_, code, err := h.Repo.Find(input.ID)
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}
	}
	email, ok := input.User["email"].(string)
	if ok {
		code, err = h.Repo.IsEmailExists(email)
		if err != nil {
			openlog.Error(err.Error())
			return common.Response{Msg: err.Error(), Data: nil, Status: code}
		}
	}
	res, code, err := h.Repo.FindAndUpdate(input.ID, input.User)
	if err != nil {
		openlog.Error(err.Error())
		return common.Response{Msg: err.Error(), Data: nil, Status: code}
	}
	return common.Response{Msg: "User Updated successfully", Data: res, Status: 201}
}
