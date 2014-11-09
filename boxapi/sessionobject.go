package boxapi

import "fmt"

type SessionObject struct {
	Type      string `json:"type"`
	Id        string `json:"id"`
	ExpiresAt string `json:"expires_at"`
}

func (so *SessionObject) IsExpired() bool {

	fmt.Printf("%s", so.ExpiresAt)
	return true
}
