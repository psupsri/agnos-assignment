package handlers

import (
	"agnos-assignment/models"
	"agnos-assignment/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PatientHandler struct {
	DB *gorm.DB
}

func filterBy(db *gorm.DB, condition bool, query any, args ...any) *gorm.DB {
	if condition {
		return db.Where(query, args...)
	}
	return db
}

func (h PatientHandler) SearchPatient(ctx *gin.Context) {
	var query models.GetPatientQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var patients []models.Patient

	db := h.DB.Model(&models.Patient{})

	if hospitalID, exists := ctx.Get("hospital_id"); exists {
		db = filterBy(db, true, "hospital_id = ?", hospitalID)
	}

	db = filterBy(db, query.NationalID != "", "nation_id = ?", query.NationalID)
	db = filterBy(db, query.PassportID != "", "passport_id = ?", query.PassportID) // ค้นหาจากชื่อจริง (ไทย หรือ อังกฤษ)
	if query.FirstName != "" {
		keyword := "%" + query.FirstName + "%"
		db = filterBy(db, true, "first_name_th LIKE ? OR first_name_en ILIKE ?", keyword, keyword)
	}

	if query.MiddleName != "" {
		keyword := "%" + query.MiddleName + "%"
		db = filterBy(db, true, "middle_name_th LIKE ? OR middle_name_en ILIKE ?", keyword, keyword)
	}

	if query.LastName != "" {
		keyword := "%" + query.LastName + "%"
		db = filterBy(db, true, "last_name_th LIKE ? OR last_name_en ILIKE ?", keyword, keyword)
	}
	db = filterBy(db, query.DateOfBirth != nil && !query.DateOfBirth.IsZero(), "date_of_birth  = ?", query.DateOfBirth)
	db = filterBy(db, query.PhoneNumber != "", "phone_number = ?", query.PhoneNumber)
	db = filterBy(db, query.Email != "", "email  = ?", query.Email)

	result := db.Find(&patients)
	if result.Error != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, result.Error.Error())
		return
	}

	utils.SuccessResponse(ctx, gin.H{
		"patients": patients,
	})
}

func (h PatientHandler) CreatePatient(ctx *gin.Context) {
	var req models.CreatePatient
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	hospitalID, ok := utils.GetHospitalID(ctx)
	if !ok {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "ไม่พบข้อมูลโรงพยาบาล")
		return
	}

	patient := models.Patient{
		PatientHN:    req.PatientHN,
		NationalID:   req.NationalID,
		PassportID:   req.PassportID,
		FirstNameTH:  req.FirstNameTH,
		MiddleNameTH: req.MiddleNameTH,
		LastNameTH:   req.LastNameTH,
		FirstNameEN:  req.FirstNameEN,
		MiddleNameEN: req.MiddleNameEN,
		LastNameEN:   req.LastNameEN,
		DateOfBirth:  req.DateOfBirth,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		Gender:       req.Gender,
		HospitalID:   hospitalID,
	}
	result := h.DB.Create(&patient)
	if result.Error != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, result.Error.Error())
		return
	}

	utils.SuccessResponse(ctx, gin.H{
		"patient": patient,
	})
}
