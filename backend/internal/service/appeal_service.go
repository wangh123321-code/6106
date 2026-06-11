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

type AppealService struct {
	appealCol     *mongo.Collection
	matchCol      *mongo.Collection
	eventCol      *mongo.Collection
	tournamentCol *mongo.Collection
	regCol        *mongo.Collection
	matchSvc      *MatchService
}

func NewAppealService(db *mongo.Database, matchSvc *MatchService) *AppealService {
	return &AppealService{
		appealCol:     db.Collection("appeals"),
		matchCol:      db.Collection("matches"),
		eventCol:      db.Collection("events"),
		tournamentCol: db.Collection("tournaments"),
		regCol:        db.Collection("registrations"),
		matchSvc:      matchSvc,
	}
}

type CreateAppealInput struct {
	MatchID  string `json:"match_id" binding:"required"`
	Reason   string `json:"reason" binding:"required"`
	Evidence string `json:"evidence" binding:"required"`
}

type ReviewAppealInput struct {
	Decision   string             `json:"decision" binding:"required,oneof=upheld rejected"`
	ReviewNote string             `json:"review_note"`
	Games      []model.GameScore  `json:"games"`
}

func (s *AppealService) Create(ctx context.Context, appellantID string, input CreateAppealInput) (*model.Appeal, error) {
	var match model.Match
	err := s.matchCol.FindOne(ctx, bson.M{"_id": input.MatchID}).Decode(&match)
	if err != nil {
		return nil, fmt.Errorf("比赛不存在")
	}

	var event model.Event
	err = s.eventCol.FindOne(ctx, bson.M{"_id": match.EventID}).Decode(&event)
	if err != nil {
		return nil, fmt.Errorf("赛事项目不存在")
	}

	var tournament model.Tournament
	err = s.tournamentCol.FindOne(ctx, bson.M{"_id": event.TournamentID}).Decode(&tournament)
	if err != nil {
		return nil, fmt.Errorf("赛事不存在")
	}

	if tournament.Status != "published" {
		return nil, fmt.Errorf("赛事成绩尚未发布，无法申诉")
	}

	if tournament.PublishedAt.IsZero() {
		return nil, fmt.Errorf("赛事发布时间无效")
	}

	if time.Since(tournament.PublishedAt) > 48*time.Hour {
		return nil, fmt.Errorf("申诉期限已过（成绩发布后48小时内可申诉）")
	}

	isPlayer := match.Player1ID == appellantID || match.Player2ID == appellantID
	if !isPlayer {
		return nil, fmt.Errorf("您不是该场比赛的参赛选手，无法申诉")
	}

	existing, _ := s.appealCol.CountDocuments(ctx, bson.M{
		"match_id":     input.MatchID,
		"appellant_id": appellantID,
		"status":       bson.M{"$in": []string{"pending", "reviewing"}},
	})
	if existing > 0 {
		return nil, fmt.Errorf("您已针对该比赛提交过申诉，正在处理中")
	}

	now := time.Now()
	appeal := model.Appeal{
		ID:          primitive.NewObjectID().Hex(),
		AppealID:    fmt.Sprintf("APL%s", primitive.NewObjectID().Hex()),
		MatchID:     input.MatchID,
		AppellantID: appellantID,
		Reason:      input.Reason,
		Evidence:    input.Evidence,
		Status:      "pending",
		CreatedAt:   now,
	}

	_, err = s.appealCol.InsertOne(ctx, appeal)
	if err != nil {
		return nil, fmt.Errorf("提交申诉失败: %v", err)
	}

	return &appeal, nil
}

func (s *AppealService) List(ctx context.Context, filter map[string]string) ([]model.Appeal, error) {
	q := bson.M{}
	if matchID, ok := filter["match_id"]; ok && matchID != "" {
		q["match_id"] = matchID
	}
	if appellantID, ok := filter["appellant_id"]; ok && appellantID != "" {
		q["appellant_id"] = appellantID
	}
	if status, ok := filter["status"]; ok && status != "" {
		q["status"] = status
	}

	cursor, err := s.appealCol.Find(ctx, q)
	if err != nil {
		return nil, err
	}
	var list []model.Appeal
	cursor.All(ctx, &list)
	return list, nil
}

func (s *AppealService) GetByID(ctx context.Context, id string) (*model.Appeal, error) {
	var a model.Appeal
	err := s.appealCol.FindOne(ctx, bson.M{"_id": id}).Decode(&a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *AppealService) Review(ctx context.Context, appealID, reviewerID, reviewerName string, input ReviewAppealInput) (*model.Appeal, error) {
	appeal, err := s.GetByID(ctx, appealID)
	if err != nil {
		return nil, fmt.Errorf("申诉不存在")
	}

	if appeal.Status == "upheld" || appeal.Status == "rejected" {
		return nil, fmt.Errorf("该申诉已处理完成")
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":       input.Decision,
			"reviewer_id":  reviewerID,
			"reviewed_at":  now,
			"review_note":  input.ReviewNote,
		},
	}

	if input.Decision == "upheld" {
		if len(input.Games) == 0 {
			return nil, fmt.Errorf("改判时必须重新录入比分")
		}
		scoreInput := ScoreInput{
			MatchID: appeal.MatchID,
			Games:   input.Games,
		}
		_, err := s.matchSvc.OverrideScore(ctx, reviewerID, reviewerName, appeal.MatchID, scoreInput)
		if err != nil {
			return nil, fmt.Errorf("改判失败: %v", err)
		}
	}

	_, err = s.appealCol.UpdateByID(ctx, appealID, update)
	if err != nil {
		return nil, fmt.Errorf("更新申诉状态失败: %v", err)
	}

	appeal.Status = input.Decision
	appeal.ReviewerID = reviewerID
	appeal.ReviewedAt = now
	appeal.ReviewNote = input.ReviewNote
	return appeal, nil
}

func (s *AppealService) HasPendingAppealsForTournament(ctx context.Context, tournamentID string) (bool, error) {
	eventsCursor, err := s.eventCol.Find(ctx, bson.M{"tournament_id": tournamentID})
	if err != nil {
		return false, err
	}
	var events []model.Event
	eventsCursor.All(ctx, &events)

	eventIDs := make([]string, 0, len(events))
	for _, e := range events {
		eventIDs = append(eventIDs, e.ID)
	}

	if len(eventIDs) == 0 {
		return false, nil
	}

	matchesCursor, err := s.matchCol.Find(ctx, bson.M{"event_id": bson.M{"$in": eventIDs}})
	if err != nil {
		return false, err
	}
	var matches []model.Match
	matchesCursor.All(ctx, &matches)

	matchIDs := make([]string, 0, len(matches))
	for _, m := range matches {
		matchIDs = append(matchIDs, m.ID)
	}

	if len(matchIDs) == 0 {
		return false, nil
	}

	count, err := s.appealCol.CountDocuments(ctx, bson.M{
		"match_id": bson.M{"$in": matchIDs},
		"status":   bson.M{"$in": []string{"pending", "reviewing"}},
	})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
