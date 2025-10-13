package unipile

type CredentialPayload struct {
	Provider string `json:"provider"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type CredentialParam struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type CookiePayload struct {
	Provider    string `json:"provider"`
	AccessToken string `json:"access_token"`
}

type CookieParam struct {
	AccessToken string `json:"access_token"`
}
