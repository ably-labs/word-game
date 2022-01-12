package model

import (
	"github.com/duo-labs/webauthn/webauthn"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID             *uint32               `gorm:"primarykey" json:"id"`
	Name           string                `json:"name"`
	CreatedAt      time.Time             `json:"createdAt"`
	UpdatedAt      time.Time             `json:"updatedAt"`
	LastActive     time.Time             `gorm:"default:CURRENT_TIMESTAMP" json:"lastActive"`
	Credentials    datatypes.JSON        `gorm:"type:json" gob:"-" json:"-"`
	CredentialsObj []webauthn.Credential `gorm:"-" gob:"-" json:"-"`
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
	return u.CredentialsObj
}

func (u *User) Exists(db *gorm.DB) bool {
	return db.Where(u).First(u).Error == nil
}
