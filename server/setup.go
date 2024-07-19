package server

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	middleware2 "github.com/labstack/echo/v4/middleware"
	"github.com/wehmoen/logart/database"
	"github.com/wehmoen/logart/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

const (
	HeaderXProject = "X-Project"
)

type Logart struct {
	e *echo.Echo
}

func NewLogart(dbUri string) (*Logart, error) {
	e := echo.New()
	e.HidePort = true
	e.HideBanner = true

	log.Printf("Connection to database %s via %s", "logart", dbUri)
	db, err := database.Open(dbUri, "logart")

	if err != nil {
		return nil, err
	}

	e.Use(db.Inject())

	e.Use(middleware2.Logger())
	e.Use(middleware.ValidateRequest())

	l := &Logart{e: e}

	e.POST("/log", l.handleLog())
	e.POST("/get", l.findLogsInTimeRange())

	return l, nil
}

func (l *Logart) findLogsInTimeRange() echo.HandlerFunc {
	return func(c echo.Context) error {

		from := c.QueryParam("from")
		to := c.QueryParam("to")

		if from == "" || to == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		// Parse time
		fromTime, err := time.Parse(time.RFC3339, from)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		toTime, err := time.Parse(time.RFC3339, to)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		if fromTime.After(toTime) {
			return c.NoContent(http.StatusBadRequest)
		}

		db := database.GetDatabaseFromContext(c)

		if db == nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		user := database.GetUserFromContext(c)

		if user == nil {
			return c.NoContent(http.StatusBadRequest)
		}

		project := c.Request().Header.Get(HeaderXProject)

		if project == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		collection := db.DB.Collection("events")

		cursor, err := collection.Find(c.Request().Context(), map[string]interface{}{
			"user_id": user.ID,
			"project": project,
			"createdAt": map[string]interface{}{
				"$gte": primitive.NewDateTimeFromTime(fromTime),
				"$lte": primitive.NewDateTimeFromTime(toTime),
			},
		})

		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		var events []database.Event

		for {
			if !cursor.Next(c.Request().Context()) {
				break
			}

			var event database.Event
			if err := cursor.Decode(&event); err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}

			events = append(events, event)
		}

		return c.JSON(http.StatusOK, events)
	}
}

func (l *Logart) handleLog() echo.HandlerFunc {
	return func(c echo.Context) error {

		db := database.GetDatabaseFromContext(c)

		if db == nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		user := database.GetUserFromContext(c)

		var data interface{}

		if err := c.Bind(&data); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		if !l.isValidJSONData(data) {
			return c.NoContent(http.StatusBadRequest)
		}

		project := c.Request().Header.Get(HeaderXProject)

		if project == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		event := database.Event{
			UserId:    user.ID,
			Project:   project,
			CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
			Data:      data,
		}

		collection := db.DB.Collection("events")
		_, err := collection.InsertOne(c.Request().Context(), event)

		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusCreated)
	}
}

func (l *Logart) isValidJSONData(data interface{}) bool {
	bytes, err := json.Marshal(data)
	if err != nil {
		return false
	}

	// Check if it's a JSON object or array
	switch data.(type) {
	case map[string]interface{}, []interface{}:
	default:
		return false
	}

	var js interface{}
	err = json.Unmarshal(bytes, &js)
	return err == nil
}

func (l *Logart) Start() {
	log.Println("Starting server on port 8080")
	l.e.Logger.Fatal(l.e.Start("127.0.0.1:8080"))
}
