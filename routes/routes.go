package routes

import (
	"agnos-assignment/handlers"
	"agnos-assignment/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUpRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	authHandler := handlers.AuthHandler{DB: db}
	staffHandler := handlers.StaffHandler{DB: db}
	hospitalHandler := handlers.HospitalHandler{DB: db}
	patientHandler := handlers.PatientHandler{DB: db}

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "ready",
		})
	})

	r.POST("/staff/login", authHandler.Login)

	r.POST("/staff/create", staffHandler.CreateStaff)

	r.GET("/staff", middleware.Auth(), staffHandler.GetAllStaff)

	r.POST("/patient", middleware.Auth(), patientHandler.CreatePatient)

	r.GET("/patient/search", middleware.Auth(), patientHandler.SearchPatient)

	r.POST("/hospital/create", middleware.Auth(), hospitalHandler.CreateHospital)

	r.GET("/hospital", middleware.Auth(), hospitalHandler.GetAllHospital)

	return r
}
