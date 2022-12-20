package models

import (
	"encoding/json"
	"time"
)

type SessionData struct {
	UserId       string
	AccessToken  string
	RefreshToken string
	Role         Role
	Exp          time.Duration
}

func (i SessionData) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}
