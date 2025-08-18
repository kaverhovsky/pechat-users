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

type UserView struct {
	ID        uuid.UUID
	NickName  string
	Firstname string
	Lastname  string
}

type UserUpdateOpts struct {
	Nickname  *string
	Firstname *string
	Lastname  *string
	Bio       *string
}
