package service

import (
	"context"
	"time"

	"github.com/tt-tournament/backend/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventService struct {
	col *mongo.Collection
}

func NewEventService(db *mongo.Database) *EventService {
	return &EventService{col: db.Collection("events")}
}

func (s *EventService) Create(ctx context.Context, tournamentID, name, eventType, drawMethod string, seedCount, bracketSize int) (*model.Event, error) {
	e := model.Event{
		ID:              primitive.NewObjectID().Hex(),
		TournamentID:    tournamentID,
		Name:            name,
		Type:            eventType,
		DrawMethod:      drawMethod,
		SeedCount:       seedCount,
		BracketSize:     bracketSize,
		Status:          "open",
		RegistrationOpen: true,
	}
	_, err := s.col.InsertOne(ctx, e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *EventService) GetByID(ctx context.Context, id string) (*model.Event, error) {
	var e model.Event
	err := s.col.FindOne(ctx, bson.M{"_id": id}).Decode(&e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *EventService) ListByTournament(ctx context.Context, tournamentID string) ([]model.Event, error) {
	cursor, err := s.col.Find(ctx, bson.M{"tournament_id": tournamentID})
	if err != nil {
		return nil, err
	}
	var list []model.Event
	cursor.All(ctx, &list)
	return list, nil
}

func (s *EventService) CloseRegistration(ctx context.Context, id string) error {
	_, err := s.col.UpdateByID(ctx, id, bson.M{"$set": bson.M{"registration_open": false, "status": "registration_closed"}})
	return err
}

func (s *EventService) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := s.col.UpdateByID(ctx, id, bson.M{"$set": bson.M{"status": status}})
	return err
}
