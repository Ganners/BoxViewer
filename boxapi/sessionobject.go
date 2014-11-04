package boxapi

type SessionObject struct {
	Type      string `json:"type"`
	Id        string `json:"id"`
	ExpiresAt string `json:"expires_at"`
}
