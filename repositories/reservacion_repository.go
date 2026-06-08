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

type ReservacionRepository struct {
	collection *mongo.Collection
}

func NewReservacionRepository() *ReservacionRepository {
	return &ReservacionRepository{collection: config.GetCollection("reservaciones")}
}

func (r *ReservacionRepository) Create(res *models.Reservacion) (*models.Reservacion, error) {
	res.ID = primitive.NewObjectID()
	res.Estado = "confirmada"
	res.CreadoEn = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *ReservacionRepository) FindByID(id string) (*models.Reservacion, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var res models.Reservacion
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *ReservacionRepository) FindByUsuario(usuarioID string, page, limit int64) ([]*models.Reservacion, error) {
	objID, err := primitive.ObjectIDFromHex(usuarioID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	skip := (page - 1) * limit
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "creado_en", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"usuario_id": objID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var reservaciones []*models.Reservacion
	if err := cursor.All(ctx, &reservaciones); err != nil {
		return nil, err
	}
	return reservaciones, nil
}

func (r *ReservacionRepository) Cancelar(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"estado": "cancelada"}},
	)
	return err
}

func (r *ReservacionRepository) CountByFuncion(funcionID string) (int64, error) {
	objID, err := primitive.ObjectIDFromHex(funcionID)
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.collection.CountDocuments(ctx, bson.M{"funcion_id": objID, "estado": "confirmada"})
}

// ReporteIngresos - aggregation: ingresos totales por función con detalle de película
func (r *ReservacionRepository) ReporteIngresos() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "estado", Value: "confirmada"}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "funciones"},
			{Key: "localField", Value: "funcion_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "funcion"},
		}}},
		{{Key: "$unwind", Value: "$funcion"}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "peliculas"},
			{Key: "localField", Value: "funcion.pelicula_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "pelicula"},
		}}},
		{{Key: "$unwind", Value: "$pelicula"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$pelicula.titulo"},
			{Key: "total_reservaciones", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "ingresos_totales", Value: bson.D{{Key: "$sum", Value: "$total"}}},
			{Key: "total_boletos", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$size", Value: "$asientos"}}}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "ingresos_totales", Value: -1}}}},
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

// ReservacionesConDetalle - lookup completo: reservacion + funcion + pelicula + usuario
func (r *ReservacionRepository) FindByIDConDetalle(id string) (bson.M, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "_id", Value: objID}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "funciones"},
			{Key: "localField", Value: "funcion_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "funcion"},
		}}},
		{{Key: "$unwind", Value: "$funcion"}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "peliculas"},
			{Key: "localField", Value: "funcion.pelicula_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "pelicula"},
		}}},
		{{Key: "$unwind", Value: "$pelicula"}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "salas"},
			{Key: "localField", Value: "funcion.sala_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "sala"},
		}}},
		{{Key: "$unwind", Value: "$sala"}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "usuarios"},
			{Key: "localField", Value: "usuario_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "usuario"},
		}}},
		{{Key: "$unwind", Value: "$usuario"}},
		{{Key: "$project", Value: bson.D{
			{Key: "usuario.password", Value: 0},
		}}},
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
	if len(results) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return results[0], nil
}
