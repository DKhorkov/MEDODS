package entities

import "time"

type RefreshToken struct {
	ID        int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	GUID      string    `json:"GUID" gorm:"unique;not null"`
	TTL       time.Time `json:"TTL" gorm:"not null"`
	Value     string    `json:"value" gorm:"unique; not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"not null"`
	DeletedAt time.Time `json:"deletedAt" gorm:"not null"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type CreateTokensDTO struct {
	GUID string `json:"GUID"`
	IP   string `json:"ip"`
}

type RefreshTokensDTO struct {
	Tokens Tokens
	IP     string `json:"ip"`
}
