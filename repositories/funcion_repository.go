package repositories

import (
	"cine-api/config"
	"cine-api/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FuncionRepository struct {
	collection *mongo.Collection
}

func NewFuncionRepository() *FuncionRepository {
	return &FuncionRepository{collection: config.GetCollection("funciones")}
}

func (r *FuncionRepository) Create(f *models.Funcion) (*models.Funcion, error) {
	f.ID = primitive.NewObjectID()
	f.Activa = true
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (r *FuncionRepository) FindByID(id string) (*models.Funcion, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var f models.Funcion
	err = r.collection.FindOne(ctx, bson.M{"_id": objID, "activa": true}).Decode(&f)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FuncionRepository) FindAll(page, limit int64) ([]*models.Funcion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	skip := (page - 1) * limit
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "fecha", Value: 1}})
	cursor, err := r.collection.Find(ctx, bson.M{"activa": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var funciones []*models.Funcion
	if err := cursor.All(ctx, &funciones); err != nil {
		return nil, err
	}
	return funciones, nil
}

func (r *FuncionRepository) FindByPelicula(peliculaID string) ([]*models.Funcion, error) {
	objID, err := primitive.ObjectIDFromHex(peliculaID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{
		"pelicula_id": objID,
		"activa":      true,
		"fecha":       bson.M{"$gte": time.Now()},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var funciones []*models.Funcion
	if err := cursor.All(ctx, &funciones); err != nil {
		return nil, err
	}
	return funciones, nil
}

func (r *FuncionRepository) FindByFecha(desde, hasta time.Time) ([]*models.Funcion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{
		"activa": true,
		"fecha":  bson.M{"$gte": desde, "$lte": hasta},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var funciones []*models.Funcion
	if err := cursor.All(ctx, &funciones); err != nil {
		return nil, err
	}
	return funciones, nil
}

func (r *FuncionRepository) UpdateAsientos(id string, cantidad int) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$inc": bson.M{"asientos_disponibles": -cantidad}},
	)
	return err
}

func (r *FuncionRepository) Delete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"activa": false}})
	return err
}

// FuncionesPorSala - aggregation: funciones agrupadas por sala con lookup
func (r *FuncionRepository) FuncionesPorSala() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "activa", Value: true}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "salas"},
			{Key: "localField", Value: "sala_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "sala"},
		}}},
		{{Key: "$unwind", Value: "$sala"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$sala.nombre"},
			{Key: "total_funciones", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "ingreso_total", Value: bson.D{{Key: "$sum", Value: "$precio"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "total_funciones", Value: -1}}}},
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
