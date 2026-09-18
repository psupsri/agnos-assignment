package models

import "time"

type Hospital struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	NameTH    string    `gorm:"not null" json:"name_th"`
	NameEn    *string   `json:"name_en"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

type CreateHospitalRequest struct {
	NameTH string  `json:"name_th" binding:"required"`
	NameEN *string `json:"name_en"`
}
