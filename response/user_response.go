package response

type UserResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Data    *UserResponseData `json:"data"`
}

type UserResponseData struct {
	User interface{} `json:"user"`
}
