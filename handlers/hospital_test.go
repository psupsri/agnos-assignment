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

func TestCreateHospital_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := utils.NewMockDB(t)
	handler := &handlers.HospitalHandler{DB: db}

	nameEnVal := "A Hospital"
	reqBody := models.CreateHospitalRequest{
		NameTH: "โรงพยาบาลเอ",
		NameEN: &nameEnVal,
	}
	jsonBytes, _ := json.Marshal(reqBody)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "hospitals"`)).
		WithArgs(
			"โรงพยาบาลเอ",
			nameEnVal,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/hospitals", handler.CreateHospital)

	req, _ := http.NewRequest(http.MethodPost, "/hospitals", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateHospital_BadRequest_MissingNameTH(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := utils.NewMockDB(t)
	handler := &handlers.HospitalHandler{DB: db}

	nameEnVal := "A Hospital"
	reqBody := models.CreateHospitalRequest{
		NameTH: "",
		NameEN: &nameEnVal,
	}
	jsonBytes, _ := json.Marshal(reqBody)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "hospitals"`)).
		WithArgs(
			"",
			nameEnVal,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/hospitals", handler.CreateHospital)

	req, _ := http.NewRequest(http.MethodPost, "/hospitals", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateHospital_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := utils.NewMockDB(t)
	handler := &handlers.HospitalHandler{DB: db}

	nameEnVal := "A Hospital"
	reqBody := models.CreateHospitalRequest{
		NameTH: "โรงพยาบาลเอ",
		NameEN: &nameEnVal,
	}
	jsonBytes, _ := json.Marshal(reqBody)

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "hospitals"`)).
		WithArgs(
			"โรงพยาบาลเอ",
			nameEnVal,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(errors.New("database insert failed"))

	mock.ExpectRollback()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/hospitals", handler.CreateHospital)

	req, _ := http.NewRequest(http.MethodPost, "/hospitals", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllHospital_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := utils.NewMockDB(t)
	handler := &handlers.HospitalHandler{DB: db}

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name_th", "name_en", "created_at", "updated_at"}).
		AddRow(1, "โรงพยาบาลเอ", "A Hospital", now, now).
		AddRow(2, "โรงพยาบาลบี", nil, now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "hospitals"`)).
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.GET("/hospitals", handler.GetAllHospital)

	req, _ := http.NewRequest(http.MethodGet, "/hospitals", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllHospital_DatabaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := utils.NewMockDB(t)
	handler := &handlers.HospitalHandler{DB: db}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "hospitals"`)).
		WillReturnError(errors.New("database connection failed"))

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.GET("/hospitals", handler.GetAllHospital)

	req, _ := http.NewRequest(http.MethodGet, "/hospitals", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
