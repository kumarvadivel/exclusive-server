package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"_id",omitempty,bson:"_id"`
	Email        string             `json:"email",bson:"email"`
	Username     string             `json:"username",bson:"username"`
	Password     string             `json:"password",bson:"password"`
	UserRole     string             `json:"userRole,omitempty",bson:"userRole" `
	UUid         string             `json:"userUuid"`
	Access_token string             `json:"access_token"`
	Firstname    string             `json:"firstname"`
	Lastname     string             `json:"lastname"`
	MiddleName   string             `json:"middlename"`
	Title        string             `json:"title"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	Mobilenumber int64              `json:"mobilenumber",bson:"mobilenumber"`
}

type UserObject struct {
	Email           string `json:"email",bson:"email"`
	Username        string `json:"username",bson:"username"`
	Password        string `json:"password",bson:"password"`
	UserRole        string `json:"userRole,omitempty",bson:"userRole" `
	UUid            string `json:"userUuid"`
	Access_token    string `json:"access_token"`
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	MiddleName      string `json:"middlename"`
	Title           string `json:"title"`
	UsernameOrEmail string `json:"usernameoremail"`
	Mobilenumber    int64  `json:"mobilenumber",bson:"mobilenumber"`
}
