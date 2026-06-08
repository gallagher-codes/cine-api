package handlers

import (
	"cine-api/models"
	"cine-api/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FuncionHandler struct {
	svc *services.FuncionService
}

func NewFuncionHandler() *FuncionHandler {
	return &FuncionHandler{svc: services.NewFuncionService()}
}

// POST /funciones
func (h *FuncionHandler) Crear(c *gin.Context) {
	var f models.Funcion
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	creada, err := h.svc.Crear(&f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, creada)
}

// GET /funciones
func (h *FuncionHandler) Listar(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	funciones, err := h.svc.ListarTodas(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": funciones, "page": page, "limit": limit})
}

// GET /funciones/:id
func (h *FuncionHandler) ObtenerPorID(c *gin.Context) {
	f, err := h.svc.ObtenerPorID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "función no encontrada"})
		return
	}
	c.JSON(http.StatusOK, f)
}

// GET /funciones/pelicula/:id
func (h *FuncionHandler) ListarPorPelicula(c *gin.Context) {
	funciones, err := h.svc.ListarPorPelicula(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": funciones})
}

// GET /funciones/fecha?desde=2025-06-01&hasta=2025-06-30
func (h *FuncionHandler) ListarPorFecha(c *gin.Context) {
	desdeStr := c.Query("desde")
	hastaStr := c.Query("hasta")
	layout := "2006-01-02"
	desde, err := time.Parse(layout, desdeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato de fecha inválido (YYYY-MM-DD)"})
		return
	}
	hasta, err := time.Parse(layout, hastaStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato de fecha inválido (YYYY-MM-DD)"})
		return
	}
	funciones, err := h.svc.ListarPorFecha(desde, hasta.Add(24*time.Hour-time.Second))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": funciones})
}

// DELETE /funciones/:id
func (h *FuncionHandler) Eliminar(c *gin.Context) {
	if err := h.svc.Eliminar(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensaje": "función eliminada correctamente"})
}

// GET /funciones/reporte/salas
func (h *FuncionHandler) ReportePorSala(c *gin.Context) {
	reporte, err := h.svc.ReportePorSala()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": reporte})
}
