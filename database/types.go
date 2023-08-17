package database

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	ApiKey  string             `bson:"api_key"`
	Enabled bool               `bson:"enabled"`
	Name    string             `bson:"name"`
}

func GetUserFromContext(c echo.Context) *User {
	user, ok := c.Get("user").(*User)
	if !ok {
		return nil
	}
	return user
}

type Event struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserId    primitive.ObjectID `bson:"user_id"`
	Project   string             `bson:"project"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	Data      interface{}        `bson:"data"`
}
