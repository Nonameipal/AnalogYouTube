package controller

type CommonResponse struct {
	Message string `json:"message"`
}

type CommonError struct {
	Error string `json:"error"`
}
