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
)

type BracketService struct {
	bracketCol *mongo.Collection
	matchCol   *mongo.Collection
	regSvc     *RegistrationService
	eventSvc   *EventService
	userSvc    *UserService
}

func NewBracketService(db *mongo.Database, regSvc *RegistrationService, eventSvc *EventService, userSvc *UserService) *BracketService {
	return &BracketService{
		bracketCol: db.Collection("brackets"),
		matchCol:   db.Collection("matches"),
		regSvc:     regSvc,
		eventSvc:   eventSvc,
		userSvc:    userSvc,
	}
}

func (s *BracketService) Generate(ctx context.Context, eventID string) (*model.Bracket, error) {
	event, err := s.eventSvc.GetByID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("项目不存在")
	}

	if event.Status != "registration_closed" {
		return nil, fmt.Errorf("报名尚未截止，无法抽签")
	}

	count, _ := s.bracketCol.CountDocuments(ctx, bson.M{"event_id": eventID})
	if count > 0 {
		return nil, fmt.Errorf("该项目已生成对阵表")
	}

	players, err := s.regSvc.GetApprovedWithRanking(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if len(players) < 2 {
		return nil, fmt.Errorf("审核通过的选手不足2人")
	}

	bracketSize := algorithm.NextPowerOf2(len(players))
	if event.BracketSize > 0 && event.BracketSize > bracketSize {
		bracketSize = event.BracketSize
	}

	var roundSlots [][]string
	if event.DrawMethod == "snake" {
		roundSlots = algorithm.SnakeSeed(players, bracketSize)
	} else {
		ids := make([]string, len(players))
		for i, p := range players {
			ids[i] = p.PlayerID
		}
		roundSlots = algorithm.RandomDraw(ids, bracketSize)
	}

	numRounds := len(roundSlots)
	allMatches := make([]model.Match, 0)
	matchMap := make(map[string]string)

	for roundIdx := 0; roundIdx < numRounds; roundIdx++ {
		roundMatchCount := len(roundSlots[roundIdx]) / 2
		for pos := 0; pos < roundMatchCount; pos++ {
			slotIdx := pos * 2
			player1ID := roundSlots[roundIdx][slotIdx]
			player2ID := roundSlots[roundIdx][slotIdx+1]

			matchID := primitive.NewObjectID().Hex()
			match := model.Match{
				ID:        matchID,
				EventID:   eventID,
				Round:     roundIdx + 1,
				Position:  pos + 1,
				Player1ID: player1ID,
				Player2ID: player2ID,
				BestOf:    7,
				Status:    "pending",
			}

			if player1ID != "" {
				user, err := s.userSvc.GetByID(ctx, player1ID)
				if err == nil {
					match.Player1Name = user.DisplayName
				}
			}
			if player2ID != "" {
				user, err := s.userSvc.GetByID(ctx, player2ID)
				if err == nil {
					match.Player2Name = user.DisplayName
				}
			}

			if player1ID == "" && player2ID != "" {
				match.WinnerID = player2ID
				match.Status = "walkover"
				match.Walkover = true
				if roundIdx+1 < numRounds {
					match.NextSlot = slotIdx
				}
			} else if player2ID == "" && player1ID != "" {
				match.WinnerID = player1ID
				match.Status = "walkover"
				match.Walkover = true
				if roundIdx+1 < numRounds {
					match.NextSlot = slotIdx
				}
			}

			key := fmt.Sprintf("r%d_p%d", roundIdx, pos)
			matchMap[key] = matchID

			allMatches = append(allMatches, match)
		}
	}

	bracketID := primitive.NewObjectID().Hex()

	for i := range allMatches {
		m := &allMatches[i]
		m.BracketID = bracketID
		if m.Round < numRounds {
			nextRound := m.Round
			nextPos := (m.Position - 1) / 2
			key := fmt.Sprintf("r%d_p%d", nextRound, nextPos)
			if nextID, ok := matchMap[key]; ok {
				m.NextMatchID = nextID
			}
		}
		m.NextSlot = (m.Position - 1) % 2
	}

	docs := make([]interface{}, len(allMatches))
	for i, m := range allMatches {
		docs[i] = m
	}
	_, err = s.matchCol.InsertMany(ctx, docs)
	if err != nil {
		return nil, err
	}

	s.propagateWalkovers(ctx, allMatches)

	bracket := model.Bracket{
		ID:        bracketID,
		EventID:   eventID,
		CreatedAt: time.Now(),
	}
	_, err = s.bracketCol.InsertOne(ctx, bracket)
	if err != nil {
		return nil, err
	}

	_ = s.eventSvc.UpdateStatus(ctx, eventID, "drawn")

	return &bracket, nil
}

func (s *BracketService) propagateWalkovers(ctx context.Context, matches []model.Match) {
	for _, m := range matches {
		if m.Walkover && m.NextMatchID != "" {
			s.advanceWinner(ctx, m.NextMatchID, m.WinnerID, m.NextSlot)
		}
	}
}

func (s *BracketService) advanceWinner(ctx context.Context, nextMatchID, winnerID string, slot int) error {
	user, _ := s.userSvc.GetByID(ctx, winnerID)
	name := ""
	if user != nil {
		name = user.DisplayName
	}

	update := bson.M{}
	if slot == 0 {
		update["player1_id"] = winnerID
		update["player1_name"] = name
	} else {
		update["player2_id"] = winnerID
		update["player2_name"] = name
	}

	_, err := s.matchCol.UpdateByID(ctx, nextMatchID, bson.M{"$set": update})
	return err
}

func (s *BracketService) GetByEvent(ctx context.Context, eventID string) (*model.Bracket, []model.Match, error) {
	var bracket model.Bracket
	err := s.bracketCol.FindOne(ctx, bson.M{"event_id": eventID}).Decode(&bracket)
	if err != nil {
		return nil, nil, err
	}

	cursor, err := s.matchCol.Find(ctx, bson.M{"event_id": eventID, "bracket_id": bracket.ID})
	if err != nil {
		return nil, nil, err
	}
	var matches []model.Match
	cursor.All(ctx, &matches)
	return &bracket, matches, nil
}

func (s *BracketService) GetByBracket(ctx context.Context, bracketID string) ([]model.Match, error) {
	cursor, err := s.matchCol.Find(ctx, bson.M{"bracket_id": bracketID})
	if err != nil {
		return nil, err
	}
	var matches []model.Match
	cursor.All(ctx, &matches)
	return matches, nil
}
