package services

import (
	"cine-api/models"
	"cine-api/repositories"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

type UsuarioService struct {
	repo *repositories.UsuarioRepository
}

func NewUsuarioService() *UsuarioService {
	return &UsuarioService{repo: repositories.NewUsuarioRepository()}
}

func (s *UsuarioService) Registrar(u *models.Usuario) (*models.Usuario, error) {
	// Verificar si email ya existe
	existing, _ := s.repo.FindByEmail(u.Email)
	if existing != nil {
		return nil, errors.New("el email ya está registrado")
	}
	return s.repo.Create(u)
}

func (s *UsuarioService) ObtenerPorID(id string) (*models.Usuario, error) {
	return s.repo.FindByID(id)
}

func (s *UsuarioService) ListarTodos(page, limit int64) ([]*models.Usuario, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	return s.repo.FindAll(page, limit)
}

func (s *UsuarioService) Actualizar(id string, campos map[string]interface{}) (*models.Usuario, error) {
	updates := bson.M{}
	if nombre, ok := campos["nombre"].(string); ok {
		updates["nombre"] = nombre
	}
	if telefono, ok := campos["telefono"].(string); ok {
		updates["telefono"] = telefono
	}
	if len(updates) == 0 {
		return nil, errors.New("no se proporcionaron campos válidos para actualizar")
	}
	return s.repo.Update(id, updates)
}

func (s *UsuarioService) Eliminar(id string) error {
	return s.repo.SoftDelete(id)
}

func (s *UsuarioService) ContarUsuarios() (int64, error) {
	return s.repo.Count()
}
