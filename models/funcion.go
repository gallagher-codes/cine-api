package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Funcion - referencia a Pelicula y Sala (relación N:1 con ambas)
type Funcion struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	PeliculaID          primitive.ObjectID `bson:"pelicula_id" json:"pelicula_id" binding:"required"` // REFERENCIA
	SalaID              primitive.ObjectID `bson:"sala_id" json:"sala_id" binding:"required"`         // REFERENCIA
	Fecha               time.Time          `bson:"fecha" json:"fecha" binding:"required"`
	Precio              float64            `bson:"precio" json:"precio" binding:"required"`
	AsientosDisponibles int                `bson:"asientos_disponibles" json:"asientos_disponibles"`
	Activa              bool               `bson:"activa" json:"activa"`
}
