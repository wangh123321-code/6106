package service

import (
	"bytes"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type RankingService struct {
	matchCol   interface{ FindOne(ctx context.Context, filter interface{}) *bson.M }
	bracketCol interface{}
	eventSvc   *EventService
	userSvc    *UserService
}

func NewRankingService(eventSvc *EventService, userSvc *UserService) *RankingService {
	return &RankingService{eventSvc: eventSvc, userSvc: userSvc}
}

type RankingResult struct {
	EventID   string `json:"event_id"`
	EventName string `json:"event_name"`
	Gold      string `json:"gold"`
	Silver    string `json:"silver"`
	Bronze1   string `json:"bronze1"`
	Bronze2   string `json:"bronze2"`
}

type MedalBoardEntry struct {
	PlayerID     string `json:"player_id"`
	DisplayName  string `json:"display_name"`
	Gold         int    `json:"gold"`
	Silver       int    `json:"silver"`
	Bronze       int    `json:"bronze"`
	Total        int    `json:"total"`
}

func GenerateMedalBoard(ctx context.Context, results []RankingResult, userSvc *UserService) []MedalBoardEntry {
	medalMap := make(map[string]*MedalBoardEntry)

	for _, r := range results {
		ids := []struct {
			id    string
			medal string
		}{
			{r.Gold, "gold"},
			{r.Silver, "silver"},
			{r.Bronze1, "bronze"},
			{r.Bronze2, "bronze"},
		}
		for _, entry := range ids {
			if entry.id == "" {
				continue
			}
			if _, ok := medalMap[entry.id]; !ok {
				user, _ := userSvc.GetByID(ctx, entry.id)
				name := entry.id
				if user != nil {
					name = user.DisplayName
				}
				medalMap[entry.id] = &MedalBoardEntry{
					PlayerID:    entry.id,
					DisplayName: name,
				}
			}
			switch entry.medal {
			case "gold":
				medalMap[entry.id].Gold++
			case "silver":
				medalMap[entry.id].Silver++
			case "bronze":
				medalMap[entry.id].Bronze++
			}
			medalMap[entry.id].Total++
		}
	}

	board := make([]MedalBoardEntry, 0, len(medalMap))
	for _, v := range medalMap {
		board = append(board, *v)
	}
	return board
}

func GenerateCertificatePDF(playerName, eventName, medal, tournamentName, date string) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("%PDF-1.4\n")
	buf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")
	buf.WriteString("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")
	buf.WriteString("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] >>\nendobj\n")

	content := fmt.Sprintf(
		"BT /F1 24 Tf 100 700 Td (%s) Tj ET BT /F1 18 Tf 100 650 Td (%s - %s) Tj ET BT /F1 14 Tf 100 600 Td (%s: %s) Tj ET BT /F1 12 Tf 100 550 Td (Date: %s) Tj ET",
		tournamentName, eventName, medal, playerName, date,
	)
	buf.WriteString(fmt.Sprintf("4 0 obj\n<< /Length %d >>\nstream\n%s\nendstream\nendobj\n", len(content), content))

	buf.WriteString("xref\n0 5\n0000000000 65535 f \n")
	buf.WriteString("trailer\n<< /Size 5 /Root 1 0 R >>\nstartxref\n0\n%%EOF\n")

	return buf.Bytes(), nil
}
