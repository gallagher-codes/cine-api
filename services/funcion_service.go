package services

import (
	"cine-api/models"
	"cine-api/repositories"
	"errors"
	"time"
)

type FuncionService struct {
	repo     *repositories.FuncionRepository
	salaRepo *repositories.SalaRepository
}

func NewFuncionService() *FuncionService {
	return &FuncionService{
		repo:     repositories.NewFuncionRepository(),
		salaRepo: repositories.NewSalaRepository(),
	}
}

func (s *FuncionService) Crear(f *models.Funcion) (*models.Funcion, error) {
	if f.Fecha.Before(time.Now()) {
		return nil, errors.New("la fecha de la función no puede ser en el pasado")
	}
	// Obtener sala para inicializar asientos disponibles
	salaID := f.SalaID.Hex()
	sala, err := s.salaRepo.FindByID(salaID)
	if err != nil {
		return nil, errors.New("sala no encontrada")
	}
	f.AsientosDisponibles = sala.Capacidad
	return s.repo.Create(f)
}

func (s *FuncionService) ObtenerPorID(id string) (*models.Funcion, error) {
	return s.repo.FindByID(id)
}

func (s *FuncionService) ListarTodas(page, limit int64) ([]*models.Funcion, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	return s.repo.FindAll(page, limit)
}

func (s *FuncionService) ListarPorPelicula(peliculaID string) ([]*models.Funcion, error) {
	return s.repo.FindByPelicula(peliculaID)
}

func (s *FuncionService) ListarPorFecha(desde, hasta time.Time) ([]*models.Funcion, error) {
	return s.repo.FindByFecha(desde, hasta)
}

func (s *FuncionService) Eliminar(id string) error {
	return s.repo.Delete(id)
}

func (s *FuncionService) ReportePorSala() (interface{}, error) {
	return s.repo.FuncionesPorSala()
}
