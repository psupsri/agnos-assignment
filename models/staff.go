package models

import "time"

type Staff struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Username   string    `gorm:"unique; not null" json:"username"`
	Password   string    `gorm:"not null" json:"password"`
	HospitalID int       `gorm:"not null" json:"hospital_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateStaffRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	HospitalID int    `json:"hospital_id" binding:"required"`
}

type GetStaffResponse struct {
	ID         int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username   string `gorm:"unique; not null" json:"username"`
	HospitalID int    `gorm:"not null" json:"hospital_id"`
}
