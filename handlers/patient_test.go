package handlers_test

import (
	"agnos-assignment/handlers"
	"agnos-assignment/models"
	"agnos-assignment/utils"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreatePatient_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := utils.NewMockDB(t)
	handler := &handlers.PatientHandler{DB: db}

	nationalId := "1100100000001"
	reqBody := models.CreatePatientRequest{
		PatientHN:   "HN001",
		NationalID:  &nationalId,
		FirstNameTH: "สมชาย",
		LastNameTH:  "ใจดี",
		FirstNameEN: "Somchai",
		LastNameEN:  "Jaidee",
	}

	jsonBytes, _ := json.Marshal(reqBody)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "patients"`)).
		WithArgs(
			"HN001",          // PatientHN
			nationalId,       // NationalID
			sqlmock.AnyArg(), // PassportID (nil)
			"สมชาย",          // FirstNameTH
			sqlmock.AnyArg(), // MiddleNameTH (nil)
			"ใจดี",           // LastNameTH
			"Somchai",        // FirstNameEN
			sqlmock.AnyArg(), // MiddleNameEN (nil)
			"Jaidee",         // LastNameEN
			sqlmock.AnyArg(), // DateOfBirth (nil)
			sqlmock.AnyArg(), // PhoneNumber (nil)
			sqlmock.AnyArg(), // Email (nil)
			sqlmock.AnyArg(), // Gender (nil)
			1,                // HospitalID (ตรงกับ int จาก Context)
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/patient", func(ctx *gin.Context) {
		ctx.Set("hospital_id", 1)
		handler.CreatePatient(ctx)
	})

	req, _ := http.NewRequest(http.MethodPost, "/patient", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePatient_MissingNameTH(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := utils.NewMockDB(t)
	handler := &handlers.PatientHandler{DB: db}

	nationalId := "1100100000001"
	reqBody := models.CreatePatientRequest{
		PatientHN:   "HN001",
		NationalID:  &nationalId,
		LastNameTH:  "ใจดี",
		FirstNameEN: "Somchai",
		LastNameEN:  "Jaidee",
	}

	jsonBytes, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/patient", func(ctx *gin.Context) {
		ctx.Set("hospital_id", 1)
		handler.CreatePatient(ctx)
	})

	req, _ := http.NewRequest(http.MethodPost, "/patient", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePatient_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := utils.NewMockDB(t)
	handler := &handlers.PatientHandler{DB: db}

	nationalId := "1100100000001"
	reqBody := models.CreatePatientRequest{
		PatientHN:   "HN001",
		NationalID:  &nationalId,
		FirstNameTH: "สมชาย",
		LastNameTH:  "ใจดี",
		FirstNameEN: "Somchai",
		LastNameEN:  "Jaidee",
	}

	jsonBytes, _ := json.Marshal(reqBody)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "patients"`)).
		WithArgs(
			"HN001",
			nationalId,
			sqlmock.AnyArg(),
			"สมชาย",
			sqlmock.AnyArg(),
			"ใจดี",
			"Somchai",
			sqlmock.AnyArg(),
			"Jaidee",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			1,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(errors.New("database insert failed"))
	mock.ExpectRollback()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/patient", func(ctx *gin.Context) {
		ctx.Set("hospital_id", 1)
		handler.CreatePatient(ctx)
	})

	req, _ := http.NewRequest(http.MethodPost, "/patient", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePatient_Unauthorized_MissingHospitalID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := utils.NewMockDB(t)
	handler := &handlers.PatientHandler{DB: db}

	nationalId := "1100100000001"
	reqBody := models.CreatePatientRequest{
		PatientHN:   "HN001",
		NationalID:  &nationalId,
		FirstNameTH: "สมชาย",
		LastNameTH:  "ใจดี",
		FirstNameEN: "Somchai",
		LastNameEN:  "Jaidee",
	}

	jsonBytes, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/patient", func(ctx *gin.Context) {
		handler.CreatePatient(ctx)
	})

	req, _ := http.NewRequest(http.MethodPost, "/patient", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSearchPatient_Success_NoFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := utils.NewMockDB(t)
	handler := &handlers.PatientHandler{DB: db}

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "patient_hn", "first_name_th", "last_name_th",
		"first_name_en", "last_name_en", "hospital_id", "created_at", "updated_at",
	}).AddRow(1, "HN001", "สมชาย", "ใจดี", "Somchai", "Jaidee", 1, now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "patients" WHERE hospital_id = $1`)).
		WithArgs(1).
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	r := gin.New()
	r.GET("/patient/search", func(ctx *gin.Context) {
		ctx.Set("hospital_id", 1)
		handler.SearchPatient(ctx)
	})

	req, _ := http.NewRequest(http.MethodGet, "/patient/search", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchPatient_Success_WithFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := utils.NewMockDB(t)
	handler := &handlers.PatientHandler{DB: db}

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "patient_hn", "first_name_th", "last_name_th",
		"first_name_en", "last_name_en", "hospital_id", "created_at", "updated_at",
	}).AddRow(1, "HN001", "สมชาย", "ใจดี", "Somchai", "Jaidee", 1, now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "patients" WHERE hospital_id = $1 AND (first_name_th LIKE $2 OR first_name_en ILIKE $3) AND (middle_name_th LIKE $4 OR middle_name_en ILIKE $5) AND (last_name_th LIKE $6 OR last_name_en ILIKE $7)`)).
		WithArgs(
			1,
			"%สมชาย%", "%สมชาย%",
			"%กลาง%", "%กลาง%",
			"%ใจดี%", "%ใจดี%",
		).
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	r := gin.New()
	r.GET("/patient/search", func(ctx *gin.Context) {
		ctx.Set("hospital_id", 1)
		handler.SearchPatient(ctx)
	})

	req, _ := http.NewRequest(http.MethodGet, "/patient/search?first_name=สมชาย&middle_name=กลาง&last_name=ใจดี", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchPatient_BadRequest_InvalidQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := utils.NewMockDB(t)
	handler := &handlers.PatientHandler{DB: db}

	w := httptest.NewRecorder()
	r := gin.New()
	r.GET("/patient/search", func(ctx *gin.Context) {
		ctx.Set("hospital_id", 1)
		handler.SearchPatient(ctx)
	})

	req, _ := http.NewRequest(http.MethodGet, "/patient/search?date_of_birth=invalid-date-format", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSearchPatient_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := utils.NewMockDB(t)
	handler := &handlers.PatientHandler{DB: db}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "patients" WHERE hospital_id = $1`)).
		WithArgs(1).
		WillReturnError(errors.New("database connection failed"))

	w := httptest.NewRecorder()
	r := gin.New()
	r.GET("/patient/search", func(ctx *gin.Context) {
		ctx.Set("hospital_id", 1)
		handler.SearchPatient(ctx)
	})

	req, _ := http.NewRequest(http.MethodGet, "/patient/search", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
