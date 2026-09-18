package handlers

import (
	"agnos-assignment/models"
	"agnos-assignment/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		DB: db,
	}
}

func (h AuthHandler) Login(ctx *gin.Context) {
	var req models.LoginRequest
	var staff models.Staff
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	result := h.DB.Where("username = ?", req.Username).First(&staff)
	if result.Error != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid username")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.Password),
		[]byte(req.Password)); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "invalid password")
		return
	}

	claims := jwt.MapClaims{
		"staff_id":    staff.ID,
		"hospital_id": staff.HospitalID,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString(
		[]byte("agnos-assignment-secret-key"),
	)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(ctx, gin.H{
		"token": signedToken,
	})
}
