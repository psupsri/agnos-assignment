package handlers_test

import (
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

	"agnos-assignment/handlers"
	"agnos-assignment/models"
	"agnos-assignment/utils"
)

func TestCreateStaff_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := utils.NewMockDB(t)
	handler := &handlers.StaffHandler{DB: db}

	reqBody := models.CreateStaffRequest{
		Username:   "admin",
		Password:   "password123",
		HospitalID: 1,
	}

	jsonBytes, _ := json.Marshal(reqBody)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "staffs"`)).
		WithArgs(
			"admin",          // Username
			sqlmock.AnyArg(), // Password (เนื่องจากผ่าน bcrypt ค่าจะเปลี่ยนตลอด ใช้ AnyArg ดักไว้)
			1,                // HospitalID
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/staff", handler.CreateStaff)

	req, _ := http.NewRequest(http.MethodPost, "/staff", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateStaff_BadRequest_MissingField(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := utils.NewMockDB(t)
	handler := &handlers.StaffHandler{DB: db}

	reqBody := models.CreateStaffRequest{
		Username:   "admin",
		HospitalID: 1,
	}

	jsonBytes, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/staff", handler.CreateStaff)

	req, _ := http.NewRequest(http.MethodPost, "/staff", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateStaff_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := utils.NewMockDB(t)
	handler := &handlers.StaffHandler{DB: db}

	reqBody := models.CreateStaffRequest{
		Username:   "admin",
		Password:   "password123",
		HospitalID: 1,
	}

	jsonBytes, _ := json.Marshal(reqBody)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "staffs"`)).
		WithArgs(
			"admin",
			sqlmock.AnyArg(),
			1,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(errors.New("database insert failed"))
	mock.ExpectRollback()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/staff", handler.CreateStaff)

	req, _ := http.NewRequest(http.MethodPost, "/staff", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateStaff_BcryptError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := utils.NewMockDB(t)
	handler := &handlers.StaffHandler{DB: db}

	longPassword := string(bytes.Repeat([]byte("a"), 73))

	reqBody := models.CreateStaffRequest{
		Username:   "admin",
		Password:   longPassword,
		HospitalID: 1,
	}

	jsonBytes, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/staff", handler.CreateStaff)

	req, _ := http.NewRequest(http.MethodPost, "/staff", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetAllStaff_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := utils.NewMockDB(t)
	handler := &handlers.StaffHandler{DB: db}

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "username", "password", "hospital_id", "created_at", "updated_at",
	}).AddRow(1, "admin", "hashed_password", 1, now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "staffs"`)).
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	r := gin.New()
	r.GET("/staff", handler.GetAllStaff)

	req, _ := http.NewRequest(http.MethodGet, "/staff", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllStaff_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := utils.NewMockDB(t)
	handler := &handlers.StaffHandler{DB: db}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "staffs"`)).
		WillReturnError(errors.New("database error"))

	w := httptest.NewRecorder()
	r := gin.New()
	r.GET("/staff", handler.GetAllStaff)

	req, _ := http.NewRequest(http.MethodGet, "/staff", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
