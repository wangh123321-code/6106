package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tt-tournament/backend/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MatchService struct {
	matchCol   *mongo.Collection
	auditCol   *mongo.Collection
	walkoverCol *mongo.Collection
	userSvc    *UserService
	bracketSvc *BracketService
}

func NewMatchService(db *mongo.Database, userSvc *UserService, bracketSvc *BracketService) *MatchService {
	return &MatchService{
		matchCol:    db.Collection("matches"),
		auditCol:    db.Collection("audit_logs"),
		walkoverCol: db.Collection("walkover_records"),
		userSvc:     userSvc,
		bracketSvc:  bracketSvc,
	}
}

type ScoreInput struct {
	MatchID string             `json:"match_id"`
	Games   []model.GameScore  `json:"games"`
}

func (s *MatchService) RecordScore(ctx context.Context, refereeID string, input ScoreInput) (*model.Match, error) {
	var match model.Match
	err := s.matchCol.FindOne(ctx, bson.M{"_id": input.MatchID}).Decode(&match)
	if err != nil {
		return nil, fmt.Errorf("比赛不存在")
	}

	if match.RefereeID != "" && match.RefereeID != refereeID {
		return nil, fmt.Errorf("您不是该场次的裁判")
	}

	if match.Status == "completed" {
		return nil, fmt.Errorf("该场次已结束")
	}

	wins1, wins2 := 0, 0
	for _, g := range input.Games {
		if g.Score1 > g.Score2 {
			wins1++
		} else if g.Score2 > g.Score1 {
			wins2++
		}
	}

	neededWins := (match.BestOf / 2) + 1
	winnerID := ""
	if wins1 >= neededWins {
		winnerID = match.Player1ID
	} else if wins2 >= neededWins {
		winnerID = match.Player2ID
	}

	if winnerID == "" {
		return nil, fmt.Errorf("比赛尚未决出胜负")
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"games":      input.Games,
			"score1":     wins1,
			"score2":     wins2,
			"winner_id":  winnerID,
			"status":     "completed",
			"referee_id": refereeID,
			"finished_at": now,
		},
	}
	_, err = s.matchCol.UpdateByID(ctx, input.MatchID, update)
	if err != nil {
		return nil, err
	}

	if match.NextMatchID != "" {
		slot := match.NextSlot
		_ = s.bracketSvc.advanceWinner(ctx, match.NextMatchID, winnerID, slot)
	}

	s.checkChampion(ctx, match.EventID, match.BracketID)

	match.Games = input.Games
	match.Score1 = wins1
	match.Score2 = wins2
	match.WinnerID = winnerID
	match.Status = "completed"
	return &match, nil
}

func (s *MatchService) Walkover(ctx context.Context, matchID, playerID, reason string) (*model.Match, error) {
	var match model.Match
	err := s.matchCol.FindOne(ctx, bson.M{"_id": matchID}).Decode(&match)
	if err != nil {
		return nil, fmt.Errorf("比赛不存在")
	}

	if match.Status == "completed" {
		return nil, fmt.Errorf("该场次已结束")
	}

	winnerID := match.Player1ID
	if playerID == match.Player1ID {
		winnerID = match.Player2ID
	}

	now := time.Now()
	_, err = s.matchCol.UpdateByID(ctx, matchID, bson.M{
		"$set": bson.M{
			"winner_id":   winnerID,
			"status":      "walkover",
			"walkover":    true,
			"finished_at": now,
		},
	})
	if err != nil {
		return nil, err
	}

	s.walkoverCol.InsertOne(ctx, model.WalkoverRecord{
		ID:        fmt.Sprintf("wo_%s", primitive.NewObjectID().Hex()),
		PlayerID:  playerID,
		MatchID:   matchID,
		EventID:   match.EventID,
		Reason:    reason,
		CreatedAt: now,
	})

	if match.NextMatchID != "" {
		_ = s.bracketSvc.advanceWinner(ctx, match.NextMatchID, winnerID, match.NextSlot)
	}

	s.checkChampion(ctx, match.EventID, match.BracketID)

	match.WinnerID = winnerID
	match.Status = "walkover"
	match.Walkover = true
	return &match, nil
}

func (s *MatchService) OverrideScore(ctx context.Context, operatorID, operatorName, matchID string, input ScoreInput) (*model.Match, error) {
	var oldMatch model.Match
	err := s.matchCol.FindOne(ctx, bson.M{"_id": matchID}).Decode(&oldMatch)
	if err != nil {
		return nil, fmt.Errorf("比赛不存在")
	}

	oldJSON, _ := json.Marshal(oldMatch)

	result, err := s.RecordScore(ctx, operatorID, input)
	if err != nil {
		return nil, err
	}

	newJSON, _ := json.Marshal(result)
	s.auditCol.InsertOne(ctx, model.AuditLog{
		ID:           primitive.NewObjectID().Hex(),
		OperatorID:   operatorID,
		OperatorName: operatorName,
		Action:       "override_score",
		TargetType:   "match",
		TargetID:     matchID,
		OldValue:     string(oldJSON),
		NewValue:     string(newJSON),
		CreatedAt:    time.Now(),
	})

	return result, nil
}

func (s *MatchService) GetByID(ctx context.Context, id string) (*model.Match, error) {
	var m model.Match
	err := s.matchCol.FindOne(ctx, bson.M{"_id": id}).Decode(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *MatchService) ListByReferee(ctx context.Context, refereeID string) ([]model.Match, error) {
	cursor, err := s.matchCol.Find(ctx, bson.M{"referee_id": refereeID})
	if err != nil {
		return nil, err
	}
	var list []model.Match
	cursor.All(ctx, &list)
	return list, nil
}

func (s *MatchService) checkChampion(ctx context.Context, eventID, bracketID string) {
	if bracketID == "" {
		return
	}
	finishedMatches, _ := s.matchCol.CountDocuments(ctx, bson.M{"bracket_id": bracketID, "status": bson.M{"$in": []string{"completed", "walkover"}}})
	totalMatches, _ := s.matchCol.CountDocuments(ctx, bson.M{"bracket_id": bracketID})

	if finishedMatches == totalMatches && totalMatches > 0 {
		var final model.Match
		err := s.matchCol.FindOne(ctx, bson.M{
			"bracket_id":    bracketID,
			"next_match_id": "",
		}).Decode(&final)
		if err == nil && final.WinnerID != "" {
			s.bracketSvc.bracketCol.UpdateByID(ctx, bracketID, bson.M{
				"$set": bson.M{"champion": final.WinnerID},
			})
		}
	}
}

func (s *MatchService) AssignReferee(ctx context.Context, matchID, refereeID string) error {
	_, err := s.matchCol.UpdateByID(ctx, matchID, bson.M{
		"$set": bson.M{"referee_id": refereeID},
	})
	return err
}

func (s *MatchService) GetAuditLogs(ctx context.Context, targetType, targetID string) ([]model.AuditLog, error) {
	filter := bson.M{}
	if targetType != "" {
		filter["target_type"] = targetType
	}
	if targetID != "" {
		filter["target_id"] = targetID
	}
	cursor, err := s.auditCol.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var logs []model.AuditLog
	cursor.All(ctx, &logs)
	return logs, nil
}
