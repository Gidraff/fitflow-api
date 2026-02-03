package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type RealmAccess struct {
	Roles []string `json:"roles"`
}

type Claims struct {
	Sub               string      `json:"sub"`
	PreferredUsername string      `json:"preferred_username"`
	RealmAccess       RealmAccess `json:"realm_access"`
	AZP               string      `json:"azp"` // Authorized party (client id)
	jwt.RegisteredClaims
}

// HasRole mimics your Rust implementation
func (c *Claims) HasRole(role string) bool {
	for _, r := range c.RealmAccess.Roles {
		if r == role {
			return true
		}
	}
	return false
}
