package repository

import (
	"context"
	"time"

	"product-api/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	FindByID(ctx context.Context, id string) (*model.Product, error)
	FindAll(ctx context.Context) ([]*model.Product, error)
	Update(ctx context.Context, id string, product *model.Product) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, name string, exactMatch bool, minPrice, maxPrice float64, sortByPrice string) ([]*model.Product, error)
}

type MongoProductRepository struct {
	collection *mongo.Collection
}

func NewMongoProductRepository(db *mongo.Database) *MongoProductRepository {
	return &MongoProductRepository{
		collection: db.Collection("products"),
	}
}

func (r *MongoProductRepository) Create(ctx context.Context, product *model.Product) error {
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return err
	}

	product.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *MongoProductRepository) FindByID(ctx context.Context, id string) (*model.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var product model.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *MongoProductRepository) FindAll(ctx context.Context) ([]*model.Product, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*model.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *MongoProductRepository) Update(ctx context.Context, id string, product *model.Product) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	product.UpdatedAt = time.Now()
	update := bson.M{
		"$set": product,
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *MongoProductRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *MongoProductRepository) Search(ctx context.Context, name string, exactMatch bool, minPrice, maxPrice float64, sortByPrice string) ([]*model.Product, error) {
	filter := bson.M{}

	if name != "" {
		if exactMatch {
			filter["name"] = name
		} else {
			filter["name"] = bson.M{"$regex": name, "$options": "i"}
		}
	}

	if minPrice > 0 || maxPrice > 0 {
		priceFilter := bson.M{}
		if minPrice > 0 {
			priceFilter["$gte"] = minPrice
		}
		if maxPrice > 0 {
			priceFilter["$lte"] = maxPrice
		}
		filter["price"] = priceFilter
	}

	findOptions := options.Find()
	if sortByPrice != "" {
		direction := 1
		if sortByPrice == "desc" {
			direction = -1
		}
		findOptions.SetSort(bson.D{{Key: "price", Value: direction}})
	}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*model.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}
