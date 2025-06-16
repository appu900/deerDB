package types

import "time"

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseDir string
}

type User struct {
	ID       string    `bson:"_id,omitempty" json:"id"`
	Username string    `bson:"username" json:"username"`
	Email    string    `bson:"email" json:"email"`
	Password string    `bson:"password" json:"-"`
	Created  time.Time `bson:"created" json:"created"`
}

