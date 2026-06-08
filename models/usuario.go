package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Usuario - colección independiente
type Usuario struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Nombre   string             `bson:"nombre" json:"nombre" binding:"required"`
	Email    string             `bson:"email" json:"email" binding:"required"`
	Telefono string             `bson:"telefono" json:"telefono"`
	Password string             `bson:"password" json:"password" binding:"required"`
	Activo   bool               `bson:"activo" json:"activo"`
	CreadoEn time.Time          `bson:"creado_en" json:"creado_en"`
}
