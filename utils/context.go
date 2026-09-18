package utils

import "github.com/gin-gonic/gin"

func GetHospitalID(ctx *gin.Context) (int, bool) {
	hospitalID, exists := ctx.Get("hospital_id")
	if !exists {
		return 0, false
	}

	switch v := hospitalID.(type) {
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	case int:
		return v, true
	default:
		return 0, false
	}
}
