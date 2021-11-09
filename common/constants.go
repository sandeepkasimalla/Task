package common

type Response struct {
	Msg    string      `json:"_msg"`
	Status int         `json:"_status"`
	Data   interface{} `json:"data"`
}

type CreateUserInput struct {
	User map[string]interface{}
}
type FetchUserInput struct {
	ID string
}

type FetchAllUsersInput struct {
	Page, Size, Filters, Sort string
}
type DeleteUserInput struct {
	ID string
}
type UpdateUserInput struct {
	ID   string
	User map[string]interface{}
}
