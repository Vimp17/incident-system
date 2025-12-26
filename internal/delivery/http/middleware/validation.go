package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
    // Регистрируем кастомные валидации
    validate.RegisterValidation("latitude", validateLatitude)
    validate.RegisterValidation("longitude", validateLongitude)
}

func validateLatitude(fl validator.FieldLevel) bool {
    lat, ok := fl.Field().Interface().(float64)
    if !ok {
        return false
    }
    return lat >= -90 && lat <= 90
}

func validateLongitude(fl validator.FieldLevel) bool {
    lng, ok := fl.Field().Interface().(float64)
    if !ok {
        return false
    }
    return lng >= -180 && lng <= 180
}

func ValidateRequest(s interface{}) gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := c.ShouldBindJSON(s); err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            c.Abort()
            return
        }
        
        if err := validate.Struct(s); err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            c.Abort()
            return
        }
        
        c.Next()
    }
}