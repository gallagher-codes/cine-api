package repositories

import (
	"cine-api/config"
	"cine-api/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SalaRepository struct {
	collection *mongo.Collection
}

func NewSalaRepository() *SalaRepository {
	return &SalaRepository{collection: config.GetCollection("salas")}
}

func (r *SalaRepository) Create(s *models.Sala) (*models.Sala, error) {
	s.ID = primitive.NewObjectID()
	s.Activa = true
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *SalaRepository) FindAll() ([]*models.Sala, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{"activa": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var salas []*models.Sala
	if err := cursor.All(ctx, &salas); err != nil {
		return nil, err
	}
	return salas, nil
}

func (r *SalaRepository) FindByID(id string) (*models.Sala, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var sala models.Sala
	err = r.collection.FindOne(ctx, bson.M{"_id": objID, "activa": true}).Decode(&sala)
	if err != nil {
		return nil, err
	}
	return &sala, nil
}

func (r *SalaRepository) FindByTipo(tipo string) ([]*models.Sala, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{"tipo": tipo, "activa": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var salas []*models.Sala
	if err := cursor.All(ctx, &salas); err != nil {
		return nil, err
	}
	return salas, nil
}

func (r *SalaRepository) Delete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"activa": false}})
	return err
}
