package application

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, app *Application) error
	GetApplications(ctx context.Context) error
}

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) Repository {
	return &repository{
		collection: db.Collection("applications"),
	}
}


func (r *repository) Create(ctx context.Context, app *Application) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, app)
	return err
}

func (r *repository) GetApplications(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor,err := r.collection.Find(ctx,bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	
	var applications []Application

	for cursor.Next(ctx) {
		var app Application
		if err := cursor.Decode(&app); err != nil {
			return err
		}
		applications = append(applications, app)
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	return nil
}