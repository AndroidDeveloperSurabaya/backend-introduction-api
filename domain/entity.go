package domain

import (
	"time"
)

// NewUser Input Boundary
type NewUser struct {
	Email    string
	Password string
	FullName string
}

// User Output Boundary
type User struct {
	Email           string    `json:"email"`
	FullName        string    `json:"fullName"`
	CreatedAt       time.Time `json:"createdAt"`
	LastModifiedAt  time.Time `json:"lastModifiedAt"`
	IsActive        bool      `json:"isActive"`
	IsEmailVerified bool      `json:"isEmailVerified"`
}

type JwtToken struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expiresIn"`
}
