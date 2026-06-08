package handlers

import (
	"cine-api/models"
	"cine-api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReservacionHandler struct {
	svc *services.ReservacionService
}

func NewReservacionHandler() *ReservacionHandler {
	return &ReservacionHandler{svc: services.NewReservacionService()}
}

// POST /reservaciones
func (h *ReservacionHandler) Crear(c *gin.Context) {
	var r models.Reservacion
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	creada, err := h.svc.Crear(&r)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, creada)
}

// GET /reservaciones/:id
func (h *ReservacionHandler) ObtenerPorID(c *gin.Context) {
	r, err := h.svc.ObtenerPorID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reservación no encontrada"})
		return
	}
	c.JSON(http.StatusOK, r)
}

// GET /reservaciones/:id/detalle
func (h *ReservacionHandler) ObtenerDetalle(c *gin.Context) {
	detalle, err := h.svc.ObtenerDetalle(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reservación no encontrada"})
		return
	}
	c.JSON(http.StatusOK, detalle)
}

// GET /reservaciones/usuario/:id
func (h *ReservacionHandler) ListarPorUsuario(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	reservaciones, err := h.svc.ListarPorUsuario(c.Param("id"), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": reservaciones, "page": page, "limit": limit})
}

// PUT /reservaciones/:id/cancelar
func (h *ReservacionHandler) Cancelar(c *gin.Context) {
	if err := h.svc.Cancelar(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensaje": "reservación cancelada correctamente"})
}

// GET /reservaciones/reporte/ingresos
func (h *ReservacionHandler) ReporteIngresos(c *gin.Context) {
	reporte, err := h.svc.ReporteIngresos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": reporte})
}
