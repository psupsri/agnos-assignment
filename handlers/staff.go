package handlers

import (
	"agnos-assignment/models"
	"agnos-assignment/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type StaffHandler struct {
	DB *gorm.DB
}

func (h *StaffHandler) CreateStaff(ctx *gin.Context) {
	var req models.CreateStaffRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	newStaff := models.Staff{
		Username:   req.Username,
		Password:   string(hashedPassword),
		HospitalID: req.HospitalID,
	}

	if err := h.DB.Create(&newStaff).Error; err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response := models.GetStaffResponse{
		ID:         newStaff.ID,
		Username:   newStaff.Username,
		HospitalID: newStaff.HospitalID,
	}

	utils.SuccessResponse(ctx, gin.H{
		"staff": response,
	})
}

func (h *StaffHandler) GetAllStaff(ctx *gin.Context) {
	var staffs []models.Staff

	result := h.DB.Find(&staffs)
	if result.Error != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, result.Error.Error())
		return
	}

	response := make([]models.GetStaffResponse, len(staffs))
	for i, staff := range staffs {
		response[i] = models.GetStaffResponse{
			ID:         staff.ID,
			Username:   staff.Username,
			HospitalID: staff.HospitalID,
		}
	}

	utils.SuccessResponse(ctx, gin.H{
		"staff": response,
	})
}
