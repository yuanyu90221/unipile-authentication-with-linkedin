package unipile

// ConnectUserWithCredentialRequest - 使用 credential 連結 user 的 request
type ConnectUserWithCredentialRequest struct {
	Account  string `json:"account" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// ConnectUserWithCookieRequest - 使用 cookie 連結 user 的 request
type ConnectUserWithCookieRequest struct {
	AccessToken string `json:"access_token" validate:"required"`
}
