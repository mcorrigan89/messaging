package entities

import (
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID
	GivenName  *string
	FamilyName *string
	Email      string
}

type NewUserEntityArgs struct {
	ID         uuid.UUID
	GivenName  *string
	FamilyName *string
	Email      string
}

func NewUserEntity(args NewUserEntityArgs) *User {
	return &User{
		ID:         args.ID,
		GivenName:  args.GivenName,
		FamilyName: args.FamilyName,
		Email:      args.Email,
	}
}
