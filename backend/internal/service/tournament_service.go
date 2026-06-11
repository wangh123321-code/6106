package service

import (
	"context"
	"fmt"
	"time"

	"github.com/tt-tournament/backend/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TournamentService struct {
	col *mongo.Collection
}

func NewTournamentService(db *mongo.Database) *TournamentService {
	return &TournamentService{col: db.Collection("tournaments")}
}

func (s *TournamentService) Create(ctx context.Context, name, location string, startDate, endDate time.Time, createdBy string) (*model.Tournament, error) {
	t := model.Tournament{
		ID:        primitive.NewObjectID().Hex(),
		Name:      name,
		StartDate: startDate,
		EndDate:   endDate,
		Location:  location,
		Status:    "draft",
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
	}
	_, err := s.col.InsertOne(ctx, t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *TournamentService) GetByID(ctx context.Context, id string) (*model.Tournament, error) {
	var t model.Tournament
	err := s.col.FindOne(ctx, bson.M{"_id": id}).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *TournamentService) List(ctx context.Context) ([]model.Tournament, error) {
	cursor, err := s.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var list []model.Tournament
	cursor.All(ctx, &list)
	return list, nil
}

func (s *TournamentService) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := s.col.UpdateByID(ctx, id, bson.M{"$set": bson.M{"status": status}})
	return err
}

func (s *TournamentService) Publish(ctx context.Context, id string) error {
	t, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if t.Status != "completed" {
		return fmt.Errorf("赛事尚未完成，无法发布成绩")
	}
	_, err = s.col.UpdateByID(ctx, id, bson.M{"$set": bson.M{
		"status":       "published",
		"published_at": time.Now(),
	}})
	return err
}
