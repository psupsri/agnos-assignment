package handlers

import (
	"agnos-assignment/models"
	"agnos-assignment/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HospitalHandler struct {
	DB *gorm.DB
}

func (h *HospitalHandler) CreateHospital(ctx *gin.Context) {
	var req models.CreateHospitalRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	newHospital := models.Hospital{
		NameTH: req.NameTH,
		NameEn: req.NameEN,
	}

	if err := h.DB.Create(&newHospital).Error; err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, gin.H{
		"staff": newHospital,
	})
}

func (h *HospitalHandler) GetAllHospital(ctx *gin.Context) {
	var hospitals []models.Hospital
	result := h.DB.Find(&hospitals)
	if result.Error != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, result.Error.Error())
		return
	}

	utils.SuccessResponse(ctx, gin.H{
		"hospitals": hospitals,
	})
}
