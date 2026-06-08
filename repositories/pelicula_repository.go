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

type PeliculaRepository struct {
	collection *mongo.Collection
}

func NewPeliculaRepository() *PeliculaRepository {
	return &PeliculaRepository{collection: config.GetCollection("peliculas")}
}

func (r *PeliculaRepository) Create(p *models.Pelicula) (*models.Pelicula, error) {
	p.ID = primitive.NewObjectID()
	p.Activo = true
	p.CreadoEn = time.Now()
	p.ActualizadoEn = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PeliculaRepository) FindByID(id string) (*models.Pelicula, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var p models.Pelicula
	err = r.collection.FindOne(ctx, bson.M{"_id": objID, "activo": true}).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PeliculaRepository) FindAll(page, limit int64) ([]*models.Pelicula, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "titulo", Value: 1}})
	cursor, err := r.collection.Find(ctx, bson.M{"activo": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var peliculas []*models.Pelicula
	if err := cursor.All(ctx, &peliculas); err != nil {
		return nil, err
	}
	return peliculas, nil
}

func (r *PeliculaRepository) FindByGenero(genero string, page, limit int64) ([]*models.Pelicula, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	skip := (page - 1) * limit
	opts := options.Find().SetSkip(skip).SetLimit(limit)
	cursor, err := r.collection.Find(ctx, bson.M{"genero": genero, "activo": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var peliculas []*models.Pelicula
	if err := cursor.All(ctx, &peliculas); err != nil {
		return nil, err
	}
	return peliculas, nil
}

func (r *PeliculaRepository) Update(id string, updates bson.M) (*models.Pelicula, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	updates["actualizado_en"] = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Pelicula
	err = r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": updates},
		opts,
	).Decode(&updated)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (r *PeliculaRepository) SoftDelete(id string) error {
	_, err := r.Update(id, bson.M{"activo": false})
	return err
}

func (r *PeliculaRepository) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.collection.CountDocuments(ctx, bson.M{"activo": true})
}

// ReportePopularidad - aggregation: películas con más reservaciones
func (r *PeliculaRepository) ReportePopularidad() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "funciones"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "pelicula_id"},
			{Key: "as", Value: "funciones"},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "reservaciones"},
			{Key: "localField", Value: "funciones._id"},
			{Key: "foreignField", Value: "funcion_id"},
			{Key: "as", Value: "reservaciones"},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "titulo", Value: 1},
			{Key: "genero", Value: 1},
			{Key: "total_reservaciones", Value: bson.D{{Key: "$size", Value: "$reservaciones"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "total_reservaciones", Value: -1}}}},
		{{Key: "$limit", Value: 10}},
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
