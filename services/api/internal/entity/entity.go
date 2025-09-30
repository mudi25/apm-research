package entity

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Message struct {
	ID              int        `json:"id"`
	EventID         string     `json:"eventId"`
	BatchID         string     `json:"batchId"`
	ProjectID       int        `json:"projectId"`
	ProjectName     string     `json:"projectName"`
	ChannelID       int        `json:"channelId"`
	ChannelName     string     `json:"channelName"`
	ChannelPlatform int16      `json:"channelPlatform"`
	TemplateID      int        `json:"templateId"`
	TemplateName    string     `json:"templateName"`
	Message         string     `json:"message"`
	Destination     string     `json:"destination"`
	CreatorID       int        `json:"creatorId"`
	CreatorName     string     `json:"creatorName"`
	Status          int16      `json:"status"`
	Attempt         int        `json:"attempt"`
	Result          string     `json:"result"`
	SendAt          *time.Time `json:"sendAt,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

type ClientDo struct {
	ID            int        `json:"id"`
	Nik           string     `json:"nik"`
	StatusNasabah string     `json:"statusNasabah"`
	StatusSend    uint8      `json:"statusSend"`
	SendAt        *time.Time `json:"sendAt,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

type Profil struct {
	ID          int       `json:"id"`
	Nama        string    `json:"nama"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Alamat      string    `json:"alamat"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
