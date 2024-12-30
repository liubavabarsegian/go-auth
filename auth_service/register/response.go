package register

type registerResponse struct {
	Status  uint8  `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

type registerErrorResponse struct {
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}
