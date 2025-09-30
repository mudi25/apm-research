package model

import (
	"research-apm/services/api/internal/entity"
	"time"
)

type Profil struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement"`
	Nama        string    `gorm:"column:nama;size:100;not null"`
	Email       string    `gorm:"column:email;size:100;not null"`
	PhoneNumber string    `gorm:"column:phone_number;size:30;not null"`
	Alamat      string    `gorm:"column:alamat;size:200;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Profil) TableName() string {
	return "apm_research.profil"
}

func (p Profil) ToEntity() entity.Profil {
	return entity.Profil{
		ID:          p.ID,
		Nama:        p.Nama,
		Email:       p.Email,
		PhoneNumber: p.PhoneNumber,
		Alamat:      p.Alamat,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
