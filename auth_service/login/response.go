package login

type loginResponse struct {
	Status       uint8  `json:"status,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type loginErrorResponse struct {
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}
