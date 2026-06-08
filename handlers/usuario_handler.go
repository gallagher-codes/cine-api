package handlers

import (
	"cine-api/models"
	"cine-api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UsuarioHandler struct {
	svc *services.UsuarioService
}

func NewUsuarioHandler() *UsuarioHandler {
	return &UsuarioHandler{svc: services.NewUsuarioService()}
}

// POST /usuarios
func (h *UsuarioHandler) Registrar(c *gin.Context) {
	var u models.Usuario
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	creado, err := h.svc.Registrar(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	creado.Password = "" // no devolver password
	c.JSON(http.StatusCreated, creado)
}

// GET /usuarios
func (h *UsuarioHandler) Listar(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	usuarios, err := h.svc.ListarTodos(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	total, _ := h.svc.ContarUsuarios()
	c.JSON(http.StatusOK, gin.H{"data": usuarios, "total": total, "page": page, "limit": limit})
}

// GET /usuarios/:id
func (h *UsuarioHandler) ObtenerPorID(c *gin.Context) {
	u, err := h.svc.ObtenerPorID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
		return
	}
	u.Password = ""
	c.JSON(http.StatusOK, u)
}

// PUT /usuarios/:id
func (h *UsuarioHandler) Actualizar(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actualizado, err := h.svc.Actualizar(c.Param("id"), body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, actualizado)
}

// DELETE /usuarios/:id
func (h *UsuarioHandler) Eliminar(c *gin.Context) {
	if err := h.svc.Eliminar(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensaje": "usuario eliminado correctamente"})
}
