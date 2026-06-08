package handlers

import (
	"cine-api/models"
	"cine-api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PeliculaHandler struct {
	svc *services.PeliculaService
}

func NewPeliculaHandler() *PeliculaHandler {
	return &PeliculaHandler{svc: services.NewPeliculaService()}
}

// POST /peliculas
func (h *PeliculaHandler) Crear(c *gin.Context) {
	var p models.Pelicula
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	creada, err := h.svc.Crear(&p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, creada)
}

// GET /peliculas
func (h *PeliculaHandler) Listar(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	peliculas, err := h.svc.ListarTodas(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	total, _ := h.svc.ContarPeliculas()
	c.JSON(http.StatusOK, gin.H{
		"data":  peliculas,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GET /peliculas/:id
func (h *PeliculaHandler) ObtenerPorID(c *gin.Context) {
	p, err := h.svc.ObtenerPorID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "película no encontrada"})
		return
	}
	c.JSON(http.StatusOK, p)
}

// GET /peliculas/genero/:genero
func (h *PeliculaHandler) ListarPorGenero(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	peliculas, err := h.svc.ListarPorGenero(c.Param("genero"), page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": peliculas})
}

// PUT /peliculas/:id
func (h *PeliculaHandler) Actualizar(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actualizada, err := h.svc.Actualizar(c.Param("id"), body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, actualizada)
}

// DELETE /peliculas/:id
func (h *PeliculaHandler) Eliminar(c *gin.Context) {
	if err := h.svc.Eliminar(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensaje": "película eliminada correctamente"})
}

// GET /peliculas/reporte/popularidad
func (h *PeliculaHandler) ReportePopularidad(c *gin.Context) {
	reporte, err := h.svc.ReportePopularidad()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": reporte})
}
