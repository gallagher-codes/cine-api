package services

import (
	"cine-api/models"
	"cine-api/repositories"
	"errors"
)

type SalaService struct {
	repo *repositories.SalaRepository
}

func NewSalaService() *SalaService {
	return &SalaService{repo: repositories.NewSalaRepository()}
}

func (s *SalaService) Crear(sala *models.Sala) (*models.Sala, error) {
	if sala.Capacidad <= 0 {
		return nil, errors.New("la capacidad debe ser mayor a 0")
	}
	tiposValidos := map[string]bool{"2D": true, "3D": true, "IMAX": true, "4DX": true}
	if sala.Tipo != "" && !tiposValidos[sala.Tipo] {
		return nil, errors.New("tipo de sala inválido (2D, 3D, IMAX, 4DX)")
	}
	return s.repo.Create(sala)
}

func (s *SalaService) ListarTodas() ([]*models.Sala, error) {
	return s.repo.FindAll()
}

func (s *SalaService) ObtenerPorID(id string) (*models.Sala, error) {
	return s.repo.FindByID(id)
}

func (s *SalaService) ListarPorTipo(tipo string) ([]*models.Sala, error) {
	return s.repo.FindByTipo(tipo)
}

func (s *SalaService) Eliminar(id string) error {
	return s.repo.Delete(id)
}
