package entities

import (
	"github.com/goccy/go-json"
	"github.com/lib/pq"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username       string         `gorm:"type:varchar(100);unique"`
	Password       string         `gorm:"type:varchar(255)"`
	Avatar         *string        `gorm:"type:text"`
	OtpActive      bool           `gorm:"type:boolean;default:false"`
	OtpVerified    bool           `gorm:"type:boolean;default:false"`
	OtpSecret      *string        `gorm:"type:varchar(255)"`
	OtpBackupCodes pq.StringArray `gorm:"type:text[];not null;default:ARRAY[]::text[]"`
	IsBlocked      bool           `gorm:"default:false"`
	IsTeam         bool           `gorm:"default:false"`
}

func (u *User) HasTwoFactor() bool {
	return u.OtpActive && u.OtpVerified
}

type RefreshToken struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	User       User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Token      string    `gorm:"type:text;not null"`
	IsRevoked  bool      `gorm:"type:boolean;default:false"`
	DeviceName string    `gorm:"type:varchar(255);not null"`
}

type Fido2Login struct {
	UserID      uuid.UUID `gorm:"type:uuid;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;primaryKey" json:"user_id"`
	DisplayName string    `gorm:"type:varchar(255);not null;unique;primaryKey" json:"display_name"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Data        string    `gorm:"type:text;not null" json:"data"`
}

func (u *Fido2Login) WebAuthnID() []byte {
	return []byte(u.UserID.String() + ":" + u.DisplayName)
}

func (u *Fido2Login) WebAuthnName() string {
	return u.DisplayName
}

func (u *Fido2Login) WebAuthnDisplayName() string {
	return u.Name
}

func (u *Fido2Login) WebAuthnCredentials() []webauthn.Credential {
	cred := &webauthn.Credential{}

	if err := json.Unmarshal([]byte(u.Data), cred); err != nil {
		panic(err)
	}

	return []webauthn.Credential{*cred}
}

func (u *Fido2Login) WebAuthnIcon() string {
	return ""
}
