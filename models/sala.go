package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Sala - colección independiente referenciada por Funcion
type Sala struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Numero    int                `bson:"numero" json:"numero" binding:"required"`
	Nombre    string             `bson:"nombre" json:"nombre"`
	Capacidad int                `bson:"capacidad" json:"capacidad" binding:"required"`
	Tipo      string             `bson:"tipo" json:"tipo"` // 2D, 3D, IMAX, 4DX
	Activa    bool               `bson:"activa" json:"activa"`
}
