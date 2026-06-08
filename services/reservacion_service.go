package services

import (
	"cine-api/models"
	"cine-api/repositories"
	"errors"
)

type ReservacionService struct {
	repo        *repositories.ReservacionRepository
	funcionRepo *repositories.FuncionRepository
}

func NewReservacionService() *ReservacionService {
	return &ReservacionService{
		repo:        repositories.NewReservacionRepository(),
		funcionRepo: repositories.NewFuncionRepository(),
	}
}

func (s *ReservacionService) Crear(res *models.Reservacion) (*models.Reservacion, error) {
	if len(res.Asientos) == 0 {
		return nil, errors.New("debe seleccionar al menos un asiento")
	}

	// Verificar disponibilidad de asientos
	funcionID := res.FuncionID.Hex()
	funcion, err := s.funcionRepo.FindByID(funcionID)
	if err != nil {
		return nil, errors.New("función no encontrada")
	}
	if funcion.AsientosDisponibles < len(res.Asientos) {
		return nil, errors.New("no hay suficientes asientos disponibles")
	}

	// Calcular total
	res.Total = funcion.Precio * float64(len(res.Asientos))

	// Crear reservación
	nueva, err := s.repo.Create(res)
	if err != nil {
		return nil, err
	}

	// Actualizar asientos disponibles
	if err := s.funcionRepo.UpdateAsientos(funcionID, len(res.Asientos)); err != nil {
		return nil, err
	}

	return nueva, nil
}

func (s *ReservacionService) ObtenerPorID(id string) (*models.Reservacion, error) {
	return s.repo.FindByID(id)
}

func (s *ReservacionService) ObtenerDetalle(id string) (interface{}, error) {
	return s.repo.FindByIDConDetalle(id)
}

func (s *ReservacionService) ListarPorUsuario(usuarioID string, page, limit int64) ([]*models.Reservacion, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	return s.repo.FindByUsuario(usuarioID, page, limit)
}

func (s *ReservacionService) Cancelar(id string) error {
	res, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("reservación no encontrada")
	}
	if res.Estado == "cancelada" {
		return errors.New("la reservación ya está cancelada")
	}
	if err := s.repo.Cancelar(id); err != nil {
		return err
	}
	// Devolver asientos a la función
	return s.funcionRepo.UpdateAsientos(res.FuncionID.Hex(), -len(res.Asientos))
}

func (s *ReservacionService) ContarPorFuncion(funcionID string) (int64, error) {
	return s.repo.CountByFuncion(funcionID)
}

func (s *ReservacionService) ReporteIngresos() (interface{}, error) {
	return s.repo.ReporteIngresos()
}
