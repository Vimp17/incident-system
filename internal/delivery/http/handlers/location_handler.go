package handlers

import (
	"net/http"

	"incident-system/internal/domain/models"
	"incident-system/internal/usecase/services"
	"incident-system/pkg/errors"

	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
    service *services.IncidentService
}

func NewLocationHandler(service *services.IncidentService) *LocationHandler {
    return &LocationHandler{service: service}
}

func (h *LocationHandler) CheckLocation(c *gin.Context) {
    var req models.LocationCheckRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, errors.NewValidationError(err))
        return
    }
    
    // Валидация координат
    if req.Latitude < -90 || req.Latitude > 90 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid latitude"})
        return
    }
    
    if req.Longitude < -180 || req.Longitude > 180 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid longitude"})
        return
    }
    
    response, err := h.service.CheckLocation(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewInternalError(err))
        return
    }
    
    c.JSON(http.StatusOK, response)
}