package boxapi

type DocumentObject struct {
	Type      string `json:"type"`
	Id        string `json:"id"`
	Status    string `json:"status"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}
