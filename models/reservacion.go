package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Asiento - subdocumento embebido dentro de Reservacion
type Asiento struct {
	Fila   string `bson:"fila" json:"fila"`
	Numero int    `bson:"numero" json:"numero"`
}

// Reservacion - relación N:M entre Usuario y Funcion (materializada)
// Referencias a usuario_id y funcion_id; asientos embebidos
type Reservacion struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UsuarioID primitive.ObjectID `bson:"usuario_id" json:"usuario_id" binding:"required"` // REFERENCIA
	FuncionID primitive.ObjectID `bson:"funcion_id" json:"funcion_id" binding:"required"` // REFERENCIA
	Asientos  []Asiento          `bson:"asientos" json:"asientos" binding:"required"`     // EMBEBIDO
	Total     float64            `bson:"total" json:"total"`
	Estado    string             `bson:"estado" json:"estado"` // confirmada, cancelada, pendiente
	CreadoEn  time.Time          `bson:"creado_en" json:"creado_en"`
}
