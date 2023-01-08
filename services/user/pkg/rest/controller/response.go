package controller

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type UserResponse struct {
	Status  string           `json:"status"`
	Message string           `json:"message,omitempty"`
	Data    UserDataResponse `json:"data,omitempty"`
}

type UserDataResponse struct {
	Username       string `json:"username"`
	Firstname      string `json:"firstname"`
	ProfilePicture string `json:"profilePicture"`
}
