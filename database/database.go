package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func Open(uri string, dbName string) (*Database, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	return &Database{Client: client, DB: db}, nil
}

func (d *Database) Inject() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("database", d)
			return next(c)
		}
	}
}

func (d *Database) UserByApiKey(apiKey string) (*User, error) {
	var user User
	collection := d.DB.Collection("users")
	filter := bson.M{"api_key": apiKey}
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *Database) CreateUser(name string) (string, error) {
	user := User{
		Name:    name,
		Enabled: true,
		ApiKey:  uuid.NewString(),
	}
	collection := d.DB.Collection("users")
	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return "", err
	}

	return user.ApiKey, nil
}

func GetDatabaseFromContext(ctx echo.Context) *Database {
	if db, ok := ctx.Get("database").(*Database); ok {
		return db
	}
	return nil
}
