package model

import (
	"research-apm/services/api/internal/entity"
	"time"
)

type ClientDo struct {
	ID            int        `gorm:"column:id;primaryKey;autoIncrement"`
	NIK           string     `gorm:"column:nik;size:50;not null"`
	StatusNasabah string     `gorm:"column:status_nasabah;size:5;not null"`
	StatusSend    uint8      `gorm:"column:status_send;not null;default:0"`
	SendAt        *time.Time `gorm:"column:send_at"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName overrides the table name used by GORM
func (ClientDo) TableName() string {
	return "REFERALNASABAH.dbo.CLIENT_DO_CLBK"
}

func (c ClientDo) ToEntity() entity.ClientDo {
	return entity.ClientDo{
		ID:            c.ID,
		Nik:           c.NIK,
		StatusNasabah: c.StatusNasabah,
		StatusSend:    c.StatusSend,
		SendAt:        c.SendAt,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}
