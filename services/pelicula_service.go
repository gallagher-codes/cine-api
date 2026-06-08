package services

import (
	"cine-api/models"
	"cine-api/repositories"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

type PeliculaService struct {
	repo *repositories.PeliculaRepository
}

func NewPeliculaService() *PeliculaService {
	return &PeliculaService{repo: repositories.NewPeliculaRepository()}
}

func (s *PeliculaService) Crear(p *models.Pelicula) (*models.Pelicula, error) {
	if p.Titulo == "" {
		return nil, errors.New("el título es obligatorio")
	}
	return s.repo.Create(p)
}

func (s *PeliculaService) ObtenerPorID(id string) (*models.Pelicula, error) {
	return s.repo.FindByID(id)
}

func (s *PeliculaService) ListarTodas(page, limit int64) ([]*models.Pelicula, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	return s.repo.FindAll(page, limit)
}

func (s *PeliculaService) ListarPorGenero(genero string, page, limit int64) ([]*models.Pelicula, error) {
	if genero == "" {
		return nil, errors.New("el género es obligatorio")
	}
	return s.repo.FindByGenero(genero, page, limit)
}

func (s *PeliculaService) Actualizar(id string, campos map[string]interface{}) (*models.Pelicula, error) {
	updates := bson.M{}
	if titulo, ok := campos["titulo"].(string); ok {
		updates["titulo"] = titulo
	}
	if sinopsis, ok := campos["sinopsis"].(string); ok {
		updates["sinopsis"] = sinopsis
	}
	if director, ok := campos["director"].(string); ok {
		updates["director"] = director
	}
	if clasificacion, ok := campos["clasificacion"].(string); ok {
		updates["clasificacion"] = clasificacion
	}
	if len(updates) == 0 {
		return nil, errors.New("no se proporcionaron campos válidos para actualizar")
	}
	return s.repo.Update(id, updates)
}

func (s *PeliculaService) Eliminar(id string) error {
	return s.repo.SoftDelete(id)
}

func (s *PeliculaService) ContarPeliculas() (int64, error) {
	return s.repo.Count()
}

func (s *PeliculaService) ReportePopularidad() (interface{}, error) {
	return s.repo.ReportePopularidad()
}
