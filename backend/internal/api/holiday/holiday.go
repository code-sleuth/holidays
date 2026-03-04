package holiday

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var contextTimeout = 5 * time.Second

const nagerURL = "https://date.nager.at/api/v3/publicholidays/2026/US"

type Nagerholiday struct {
	Date      string `json:"date"`
	LocalName string `json:"localName"`
	Name      string `json:"name"`
}

type Lookup struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	HolidayName string             `json:"holidayName" bson:"holidayName"`
	HolidayDate string             `json:"holidayDate" bson:"holidayDate"`
	DaysUntil   int                `json:"daysUntil" bson:"daysUntil"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
}

type CreateLookupRequest struct {
	HolidayName string `json:"holidayName"`
	HolidayDate string `json:"holidayDate"`
}

// Handler & Service implementation

type Handler struct {
	DB           *mongo.Client
	DatabaseName string
	Logger       zerolog.Logger
}

type Service interface {
	GetHolidays(c *echo.Context) error
	CreateLookup(c *echo.Context) error
	GetLookups(c *echo.Context) error
}

func NewManager(h Handler) *Handler {
	return &Handler{
		DB:           h.DB,
		DatabaseName: h.DatabaseName,
		Logger:       h.Logger,
	}
}

func (h *Handler) lookupCollection() *mongo.Collection {
	return h.DB.Database(h.DatabaseName).Collection("lookups")
}

// GET /api/v1/holidays (proxy to Nager.Date API)
func (h *Handler) GetHolidays(c *echo.Context) error {
	resp, err := http.Get(nagerURL)
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed to fetch holidays from Nager.Date")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch holidays"})
	}
	defer resp.Body.Close()

	var holidays []Nagerholiday
	if err := json.NewDecoder(resp.Body).Decode(&holidays); err != nil {
		h.Logger.Error().Err(err).Msg("failed to decode holidays response")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to decode holidays"})
	}

	return c.JSON(http.StatusOK, holidays)
}

// POST /api/v1/lookups (calculate days until holiday)
func (h *Handler) CreateLookUp(c *echo.Context) error {
	var req CreateLookupRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	holidayDate, err := time.Parse("2006-01-02", req.HolidayDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid date format"})
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	target := time.Date(holidayDate.Year(), holidayDate.Month(), holidayDate.Day(), 0, 0, 0, 0, time.UTC)
	daysUntil := int(math.Ceil(float64(target.Sub(today).Hours() / 24)))

	ctx, cancel := context.WithTimeout(c.Request().Context(), contextTimeout)
	defer cancel()

	lookup := Lookup{
		HolidayName: req.HolidayName,
		HolidayDate: req.HolidayDate,
		DaysUntil:   daysUntil,
		CreatedAt:   time.Now(),
	}

	res, err := h.lookupCollection().InsertOne(ctx, lookup)
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed to save lookup")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save lookup"})
	}

	lookup.ID = res.InsertedID.(primitive.ObjectID)
	return c.JSON(http.StatusCreated, lookup)
}

// GET /api/v1/lookups (return all saved lookups)
func (h *Handler) GetLookups(c *echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), contextTimeout)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := h.lookupCollection().Find(ctx, bson.M{}, opts)
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed to fetch lookups")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch lookups"})
	}
	defer cursor.Close(ctx)

	lookups := []Lookup{}
	if err := cursor.All(ctx, &lookups); err != nil {
		h.Logger.Error().Err(err).Msg("failed tp decode lookups")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed tp decode lookups"})
	}

	return c.JSON(http.StatusOK, lookups)
}
