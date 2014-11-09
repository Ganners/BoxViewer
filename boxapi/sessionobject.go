package boxapi

import (
	"log"
	"time"
)

type SessionObject struct {
	Type      string `json:"type"`
	Id        string `json:"id"`
	ExpiresAt string `json:"expires_at"`
}

func (so *SessionObject) IsExpired() bool {

	t, err := time.Parse(time.RFC3339Nano, so.ExpiresAt)
	if err != nil {
		log.Fatal("Problem parsing date on ", so.ExpiresAt)
	}

	if t.After(time.Now()) {
		return false
	}

	return true
}
