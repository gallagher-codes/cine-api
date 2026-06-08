package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Actor - subdocumento embebido dentro de Pelicula
type Actor struct {
	Nombre    string `bson:"nombre" json:"nombre"`
	Personaje string `bson:"personaje" json:"personaje"`
}

// Pelicula - colección principal con reparto embebido
type Pelicula struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Titulo        string             `bson:"titulo" json:"titulo" binding:"required"`
	Genero        []string           `bson:"genero" json:"genero"`
	Duracion      int                `bson:"duracion" json:"duracion"` // minutos
	Clasificacion string             `bson:"clasificacion" json:"clasificacion"` // AA, A, B, B15, C
	Director      string             `bson:"director" json:"director"`
	Sinopsis      string             `bson:"sinopsis" json:"sinopsis"`
	Reparto       []Actor            `bson:"reparto" json:"reparto"` // EMBEBIDO
	Activo        bool               `bson:"activo" json:"activo"`
	CreadoEn      time.Time          `bson:"creado_en" json:"creado_en"`
	ActualizadoEn time.Time          `bson:"actualizado_en" json:"actualizado_en"`
}
