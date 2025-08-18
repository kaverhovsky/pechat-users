package model

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Nickname     string
	PasswordHash string
	Firstname    string
	Lastname     string
	Email        string
	Bio          string
	// avatar id
	// role (?)
	// phone number (?)
}
