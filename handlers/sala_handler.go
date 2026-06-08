package handlers

import (
	"cine-api/models"
	"cine-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SalaHandler struct {
	svc *services.SalaService
}

func NewSalaHandler() *SalaHandler {
	return &SalaHandler{svc: services.NewSalaService()}
}

// POST /salas
func (h *SalaHandler) Crear(c *gin.Context) {
	var s models.Sala
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	creada, err := h.svc.Crear(&s)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, creada)
}

// GET /salas
func (h *SalaHandler) Listar(c *gin.Context) {
	salas, err := h.svc.ListarTodas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": salas})
}

// GET /salas/:id
func (h *SalaHandler) ObtenerPorID(c *gin.Context) {
	sala, err := h.svc.ObtenerPorID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sala no encontrada"})
		return
	}
	c.JSON(http.StatusOK, sala)
}

// GET /salas/tipo/:tipo
func (h *SalaHandler) ListarPorTipo(c *gin.Context) {
	salas, err := h.svc.ListarPorTipo(c.Param("tipo"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": salas})
}

// DELETE /salas/:id
func (h *SalaHandler) Eliminar(c *gin.Context) {
	if err := h.svc.Eliminar(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensaje": "sala eliminada correctamente"})
}
