package unipile

type ConnectResult struct {
	Object     string     `json:"object"`
	AccountID  string     `json:"account_id"`
	CheckPoint CheckPoint `json:"checkpoint,omitempty"`
}

type CheckPoint struct {
	Type string `json:"type"`
}
