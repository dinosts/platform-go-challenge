package server

type ErrorResponse struct {
	Error string `json:"error"`
}

type DataResponse struct {
	Data any `json:"data"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
