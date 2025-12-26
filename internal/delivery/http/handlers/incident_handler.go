package handlers

import (
	"net/http"
	"strconv"

	"incident-system/internal/domain/models"
	"incident-system/internal/usecase/services"
	"incident-system/pkg/errors"

	"github.com/gin-gonic/gin"
)

type IncidentHandler struct {
    service *services.IncidentService
}

func NewIncidentHandler(service *services.IncidentService) *IncidentHandler {
    return &IncidentHandler{service: service}
}

func (h *IncidentHandler) CreateIncident(c *gin.Context) {
    var req models.CreateIncidentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, errors.NewValidationError(err))
        return
    }
    
    incident, err := h.service.CreateIncident(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewInternalError(err))
        return
    }
    
    c.JSON(http.StatusCreated, incident)
}

func (h *IncidentHandler) GetIncident(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, errors.NewValidationError(err))
        return
    }
    
    incident, err := h.service.GetIncident(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewInternalError(err))
        return
    }
    
    if incident == nil {
        c.JSON(http.StatusNotFound, errors.NewNotFoundError("incident"))
        return
    }
    
    c.JSON(http.StatusOK, incident)
}

func (h *IncidentHandler) ListIncidents(c *gin.Context) {
    // Параметры пагинации
    limitStr := c.DefaultQuery("limit", "10")
    pageStr := c.DefaultQuery("page", "1")
    activeOnlyStr := c.DefaultQuery("active_only", "true")
    
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit < 1 || limit > 100 {
        limit = 10
    }
    
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }
    
    activeOnly, err := strconv.ParseBool(activeOnlyStr)
    if err != nil {
        activeOnly = true
    }
    
    offset := (page - 1) * limit
    
    incidents, total, err := h.service.ListIncidents(c.Request.Context(), limit, offset, activeOnly)
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewInternalError(err))
        return
    }
    
    // Расчет метаданных пагинации
    totalPages := (total + limit - 1) / limit
    if totalPages == 0 {
        totalPages = 1
    }
    
    response := gin.H{
        "data": incidents,
        "meta": gin.H{
            "page":        page,
            "limit":       limit,
            "total":       total,
            "total_pages": totalPages,
            "has_next":    page < totalPages,
            "has_prev":    page > 1,
        },
    }
    
    c.JSON(http.StatusOK, response)
}

func (h *IncidentHandler) UpdateIncident(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, errors.NewValidationError(err))
        return
    }
    
    var req models.UpdateIncidentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, errors.NewValidationError(err))
        return
    }
    
    incident, err := h.service.UpdateIncident(c.Request.Context(), id, req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewInternalError(err))
        return
    }
    
    if incident == nil {
        c.JSON(http.StatusNotFound, errors.NewNotFoundError("incident"))
        return
    }
    
    c.JSON(http.StatusOK, incident)
}

func (h *IncidentHandler) DeleteIncident(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, errors.NewValidationError(err))
        return
    }
    
    if err := h.service.DeleteIncident(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewInternalError(err))
        return
    }
    
    c.Status(http.StatusNoContent)
}

func (h *IncidentHandler) GetStats(c *gin.Context) {
    minutesStr := c.DefaultQuery("minutes", "60")
    minutes, err := strconv.Atoi(minutesStr)
    if err != nil || minutes < 1 {
        minutes = 60
    }
    
    stats, err := h.service.GetStats(c.Request.Context(), minutes)
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewInternalError(err))
        return
    }
    
    c.JSON(http.StatusOK, stats)
}