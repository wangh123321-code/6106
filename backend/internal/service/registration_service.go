package service

import (
	"context"
	"fmt"
	"time"

	"github.com/tt-tournament/backend/internal/algorithm"
	"github.com/tt-tournament/backend/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RegistrationService struct {
	col  *mongo.Collection
	user *UserService
}

func NewRegistrationService(db *mongo.Database, userSvc *UserService) *RegistrationService {
	return &RegistrationService{col: db.Collection("registrations"), user: userSvc}
}

func (s *RegistrationService) Register(ctx context.Context, eventID, playerID, partnerID, qualificationURL string) (*model.Registration, error) {
	count, _ := s.col.CountDocuments(ctx, bson.M{"event_id": eventID, "player_id": playerID})
	if count > 0 {
		return nil, fmt.Errorf("已报名该项目")
	}

	reg := model.Registration{
		ID:               primitive.NewObjectID().Hex(),
		EventID:          eventID,
		PlayerID:         playerID,
		PartnerID:        partnerID,
		QualificationURL: qualificationURL,
		Status:           "pending",
		CreatedAt:        time.Now(),
	}
	_, err := s.col.InsertOne(ctx, reg)
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func (s *RegistrationService) BatchApprove(ctx context.Context, eventID, reviewerID string, regIDs []string) (int, error) {
	now := time.Now()
	filter := bson.M{
		"event_id": eventID,
		"_id":      bson.M{"$in": regIDs},
		"status":   "pending",
	}
	update := bson.M{
		"$set": bson.M{
			"status":      "approved",
			"reviewed_by": reviewerID,
			"reviewed_at": now,
		},
	}
	result, err := s.col.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return int(result.ModifiedCount), nil
}

func (s *RegistrationService) Reject(ctx context.Context, regID, reviewerID string) error {
	_, err := s.col.UpdateByID(ctx, regID, bson.M{
		"$set": bson.M{
			"status":      "rejected",
			"reviewed_by": reviewerID,
			"reviewed_at": time.Now(),
		},
	})
	return err
}

func (s *RegistrationService) ListByEvent(ctx context.Context, eventID string) ([]model.Registration, error) {
	opts := options.Find().SetSort(bson.M{"created_at": 1})
	cursor, err := s.col.Find(ctx, bson.M{"event_id": eventID}, opts)
	if err != nil {
		return nil, err
	}
	var list []model.Registration
	cursor.All(ctx, &list)
	return list, nil
}

func (s *RegistrationService) GetApprovedPlayerIDs(ctx context.Context, eventID string) ([]string, error) {
	cursor, err := s.col.Find(ctx, bson.M{"event_id": eventID, "status": "approved"})
	if err != nil {
		return nil, err
	}
	var regs []model.Registration
	cursor.All(ctx, &regs)

	ids := make([]string, 0, len(regs))
	for _, r := range regs {
		ids = append(ids, r.PlayerID)
		if r.PartnerID != "" {
			ids = append(ids, r.PartnerID)
		}
	}
	return ids, nil
}

func (s *RegistrationService) GetApprovedWithRanking(ctx context.Context, eventID string) ([]algorithm.SeedEntry, error) {
	cursor, err := s.col.Find(ctx, bson.M{"event_id": eventID, "status": "approved"})
	if err != nil {
		return nil, err
	}
	var regs []model.Registration
	cursor.All(ctx, &regs)

	entries := make([]algorithm.SeedEntry, 0, len(regs))
	for _, r := range regs {
		user, err := s.user.GetByID(ctx, r.PlayerID)
		ranking := 9999
		if err == nil {
			ranking = user.Ranking
			if ranking == 0 {
				ranking = 9999
			}
		}
		entries = append(entries, algorithm.SeedEntry{
			PlayerID: r.PlayerID,
			Ranking:  ranking,
		})
	}
	return entries, nil
}
