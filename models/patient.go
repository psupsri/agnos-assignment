package models

import "time"

type GenderType string

const (
	GenderMale   GenderType = "M"
	GenderFemale GenderType = "F"
)

type Patient struct {
	ID           int         `gorm:"primaryKey;autoIncrement" json:"id"`
	PatientHN    string      `gorm:"not null;unique" json:"patient_hn"`
	NationalID   *string     `gorm:"unique" json:"national_id"`
	PassportID   *string     `gorm:"unique" json:"passport_id"`
	FirstNameTH  string      `gorm:"not null" json:"first_name_th"`
	MiddleNameTH *string     `json:"middle_name_th"`
	LastNameTH   string      `gorm:"not null" json:"last_name_th"`
	FirstNameEN  string      `gorm:"not null" json:"first_name_en"`
	MiddleNameEN *string     `json:"middle_name_en"`
	LastNameEN   string      `gorm:"not null" json:"last_name_en"`
	DateOfBirth  *time.Time  `json:"date_of_birth"`
	PhoneNumber  *string     `json:"phone_number"`
	Email        *string     `json:"email"`
	Gender       *GenderType `gorm:"type:varchar(1);check:gender IN ('M', 'F')" json:"gender"`
	HospitalID   int         `gorm:"not null" json:"hospital_id"`
	CreatedAt    time.Time   `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"not null" json:"updated_at"`
}

type GetPatientQuery struct {
	NationalID  string     `form:"nation_id"`
	PassportID  string     `form:"passport_id"`
	FirstName   string     `form:"first_name"`
	MiddleName  string     `form:"middle_name"`
	LastName    string     `form:"last_name"`
	DateOfBirth *time.Time `form:"date_of_birth"`
	PhoneNumber string     `form:"phone_number"`
	Email       string     `form:"email"`
}

type CreatePatientRequest struct {
	PatientHN    string      `json:"patient_hn" binding:"required"`
	NationalID   *string     `json:"national_id"`
	PassportID   *string     `json:"passport_id"`
	FirstNameTH  string      `json:"first_name_th" binding:"required"`
	MiddleNameTH *string     `json:"middle_name_th"`
	LastNameTH   string      `json:"last_name_th" binding:"required"`
	FirstNameEN  string      `json:"first_name_en" binding:"required"`
	MiddleNameEN *string     `json:"middle_name_en"`
	LastNameEN   string      `json:"last_name_en" binding:"required"`
	DateOfBirth  *time.Time  `json:"date_of_birth"`
	PhoneNumber  *string     `json:"phone_number"`
	Email        *string     `json:"email"`
	Gender       *GenderType `json:"gender" binding:"omitempty,oneof=M F"`
}
