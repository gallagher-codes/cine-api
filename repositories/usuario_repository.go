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

type UsuarioRepository struct {
	collection *mongo.Collection
}

func NewUsuarioRepository() *UsuarioRepository {
	return &UsuarioRepository{collection: config.GetCollection("usuarios")}
}

func (r *UsuarioRepository) Create(u *models.Usuario) (*models.Usuario, error) {
	u.ID = primitive.NewObjectID()
	u.Activo = true
	u.CreadoEn = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UsuarioRepository) FindByID(id string) (*models.Usuario, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var u models.Usuario
	err = r.collection.FindOne(ctx, bson.M{"_id": objID, "activo": true}).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsuarioRepository) FindByEmail(email string) (*models.Usuario, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var u models.Usuario
	err := r.collection.FindOne(ctx, bson.M{"email": email, "activo": true}).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsuarioRepository) FindAll(page, limit int64) ([]*models.Usuario, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetProjection(bson.M{"password": 0}) // no exponer password
	cursor, err := r.collection.Find(ctx, bson.M{"activo": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var usuarios []*models.Usuario
	if err := cursor.All(ctx, &usuarios); err != nil {
		return nil, err
	}
	return usuarios, nil
}

func (r *UsuarioRepository) Update(id string, updates bson.M) (*models.Usuario, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetProjection(bson.M{"password": 0})
	var updated models.Usuario
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

func (r *UsuarioRepository) SoftDelete(id string) error {
	_, err := r.Update(id, bson.M{"activo": false})
	return err
}

func (r *UsuarioRepository) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.collection.CountDocuments(ctx, bson.M{"activo": true})
}
