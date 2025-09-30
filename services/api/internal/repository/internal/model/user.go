package model

import (
	"research-apm/services/api/internal/entity"
	"time"
)

type User struct {
	ID        string    `bson:"_id"`
	Name      string    `bson:"name"`
	Address   string    `bson:"address"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (u User) ToEntity() entity.User {
	return entity.User{
		ID:        u.ID,
		Name:      u.Name,
		Address:   u.Address,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func NewUser(u entity.User) User {
	return User{
		ID:        u.ID,
		Name:      u.Name,
		Address:   u.Address,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
