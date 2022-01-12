package model

import (
	"github.com/duo-labs/webauthn/webauthn"
	"time"
)

type User struct {
	ID          *uint32 `gorm:"primarykey"`
	Name        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LastActive  time.Time             `gorm:"default:CURRENT_TIMESTAMP"`
	Credentials []webauthn.Credential `gorm:"type:json"`
}

func (u *User) WebAuthnID() []byte {
	return []byte(u.Name)
}

func (u *User) WebAuthnName() string {
	return u.Name
}

func (u *User) WebAuthnDisplayName() string {
	return u.Name
}

func (u *User) WebAuthnIcon() string {
	return ""
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}
