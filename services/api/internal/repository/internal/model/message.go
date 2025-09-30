package model

import (
	"research-apm/services/api/internal/entity"
	"time"
)

type Message struct {
	ID              int        `gorm:"primaryKey;autoIncrement;column:id"`
	EventID         string     `gorm:"size:40;not null;column:event_id"`
	BatchID         string     `gorm:"size:40;not null;column:batch_id"`
	ProjectID       int        `gorm:"not null;column:project_id"`
	ProjectName     string     `gorm:"size:50;not null;column:project_name"`
	ChannelID       int        `gorm:"not null;column:channel_id"`
	ChannelName     string     `gorm:"size:50;not null;column:channel_name"`
	ChannelPlatform int16      `gorm:"not null;column:channel_platform"`
	TemplateID      int        `gorm:"not null;column:template_id"`
	TemplateName    string     `gorm:"size:50;not null;column:template_name"`
	Message         string     `gorm:"type:text;not null;column:message"`
	Destination     string     `gorm:"type:text;not null;column:destination"`
	CreatorID       int        `gorm:"not null;column:creator_id"`
	CreatorName     string     `gorm:"size:50;not null;column:creator_name"`
	Status          int16      `gorm:"not null;column:status"`
	Attempt         int        `gorm:"not null;column:attempt"`
	Result          string     `gorm:"type:text;not null;column:result"`
	SendAt          *time.Time `gorm:"column:send_at"`
	CreatedAt       time.Time  `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime;column:updated_at"`
}

func (Message) TableName() string {
	return "public.messages"
}

func (m Message) ToEntity() entity.Message {
	return entity.Message{
		ID:              m.ID,
		EventID:         m.EventID,
		BatchID:         m.BatchID,
		ProjectID:       m.ProjectID,
		ProjectName:     m.ProjectName,
		ChannelID:       m.ChannelID,
		ChannelName:     m.ChannelName,
		ChannelPlatform: m.ChannelPlatform,
		TemplateID:      m.TemplateID,
		TemplateName:    m.TemplateName,
		Message:         m.Message,
		Destination:     m.Destination,
		CreatorID:       m.CreatorID,
		CreatorName:     m.CreatorName,
		Status:          m.Status,
		Attempt:         m.Attempt,
		Result:          m.Result,
		SendAt:          m.SendAt,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}
