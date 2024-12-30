package logout

type logoutResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type logoutErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}
