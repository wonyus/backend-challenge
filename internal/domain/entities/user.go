package entities

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func NewUser(name, email, hashedPassword string) (*User, error) {
	if name == "" || email == "" || hashedPassword == "" {
		return nil, errors.New("name, email, and password are required")
	}

	now := time.Now()
	return &User{
		ID:        primitive.NewObjectID(),
		Name:      name,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) UpdateName(name string) {
	u.Name = name
	u.UpdatedAt = time.Now()
}

func (u *User) UpdateEmail(email string) {
	u.Email = email
	u.UpdatedAt = time.Now()
}
