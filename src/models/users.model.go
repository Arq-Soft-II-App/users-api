package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Lastname  string             `bson:"lastname"`
	Birthdate time.Time          `bson:"birthdate"`
	Role      string             `bson:"role"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Avatar    string             `bson:"avatar"`
}
