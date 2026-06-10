package model

import "time"

type User struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Username    string    `bson:"username" json:"username"`
	Password    string    `bson:"password" json:"-"`
	DisplayName string    `bson:"display_name" json:"display_name"`
	Role        string    `bson:"role" json:"role"`
	Ranking     int       `bson:"ranking" json:"ranking"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}

type Tournament struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Name        string    `bson:"name" json:"name"`
	StartDate   time.Time `bson:"start_date" json:"start_date"`
	EndDate     time.Time `bson:"end_date" json:"end_date"`
	Location    string    `bson:"location" json:"location"`
	Status      string    `bson:"status" json:"status"`
	CreatedBy   string    `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}

type Event struct {
	ID              string      `bson:"_id,omitempty" json:"id"`
	TournamentID    string      `bson:"tournament_id" json:"tournament_id"`
	Name            string      `bson:"name" json:"name"`
	Type            string      `bson:"type" json:"type"`
	DrawMethod      string      `bson:"draw_method" json:"draw_method"`
	SeedCount       int         `bson:"seed_count" json:"seed_count"`
	BracketSize     int         `bson:"bracket_size" json:"bracket_size"`
	Status          string      `bson:"status" json:"status"`
	RegistrationOpen bool        `bson:"registration_open" json:"registration_open"`
}

type Registration struct {
	ID               string    `bson:"_id,omitempty" json:"id"`
	EventID          string    `bson:"event_id" json:"event_id"`
	PlayerID         string    `bson:"player_id" json:"player_id"`
	PartnerID        string    `bson:"partner_id,omitempty" json:"partner_id,omitempty"`
	QualificationURL string    `bson:"qualification_url" json:"qualification_url"`
	Status           string    `bson:"status" json:"status"`
	ReviewedBy       string    `bson:"reviewed_by,omitempty" json:"reviewed_by,omitempty"`
	ReviewedAt       time.Time `bson:"reviewed_at,omitempty" json:"reviewed_at,omitempty"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
}

type Bracket struct {
	ID        string           `bson:"_id,omitempty" json:"id"`
	EventID   string           `bson:"event_id" json:"event_id"`
	Rounds    [][]MatchSlot    `bson:"rounds" json:"rounds"`
	Champion  *string          `bson:"champion,omitempty" json:"champion,omitempty"`
	CreatedAt time.Time        `bson:"created_at" json:"created_at"`
}

type MatchSlot struct {
	MatchID    string `bson:"match_id" json:"match_id"`
	Player1ID  string `bson:"player1_id,omitempty" json:"player1_id,omitempty"`
	Player1Name string `bson:"player1_name,omitempty" json:"player1_name,omitempty"`
	Player2ID  string `bson:"player2_id,omitempty" json:"player2_id,omitempty"`
	Player2Name string `bson:"player2_name,omitempty" json:"player2_name,omitempty"`
	Score1     *int   `bson:"score1,omitempty" json:"score1,omitempty"`
	Score2     *int   `bson:"score2,omitempty" json:"score2,omitempty"`
	WinnerID   string `bson:"winner_id,omitempty" json:"winner_id,omitempty"`
	Status     string `bson:"status" json:"status"`
	NextMatchID string `bson:"next_match_id,omitempty" json:"next_match_id,omitempty"`
}

type Match struct {
	ID           string    `bson:"_id,omitempty" json:"id"`
	BracketID    string    `bson:"bracket_id" json:"bracket_id"`
	EventID      string    `bson:"event_id" json:"event_id"`
	Round        int       `bson:"round" json:"round"`
	Position     int       `bson:"position" json:"position"`
	Player1ID    string    `bson:"player1_id,omitempty" json:"player1_id,omitempty"`
	Player1Name  string    `bson:"player1_name,omitempty" json:"player1_name,omitempty"`
	Player2ID    string    `bson:"player2_id,omitempty" json:"player2_id,omitempty"`
	Player2Name  string    `bson:"player2_name,omitempty" json:"player2_name,omitempty"`
	Score1       int       `bson:"score1" json:"score1"`
	Score2       int       `bson:"score2" json:"score2"`
	BestOf       int       `bson:"best_of" json:"best_of"`
	Games        []GameScore `bson:"games" json:"games"`
	WinnerID     string    `bson:"winner_id,omitempty" json:"winner_id,omitempty"`
	Status       string    `bson:"status" json:"status"`
	RefereeID    string    `bson:"referee_id,omitempty" json:"referee_id,omitempty"`
	NextMatchID  string    `bson:"next_match_id,omitempty" json:"next_match_id,omitempty"`
	NextSlot     int       `bson:"next_slot" json:"next_slot"`
	Walkover     bool      `bson:"walkover" json:"walkover"`
	ScheduledAt  *time.Time `bson:"scheduled_at,omitempty" json:"scheduled_at,omitempty"`
	FinishedAt   *time.Time `bson:"finished_at,omitempty" json:"finished_at,omitempty"`
}

type GameScore struct {
	Score1 int `bson:"score1" json:"score1"`
	Score2 int `bson:"score2" json:"score2"`
}

type AuditLog struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	OperatorID  string    `bson:"operator_id" json:"operator_id"`
	OperatorName string `bson:"operator_name" json:"operator_name"`
	Action      string    `bson:"action" json:"action"`
	TargetType  string    `bson:"target_type" json:"target_type"`
	TargetID    string    `bson:"target_id" json:"target_id"`
	OldValue    string    `bson:"old_value,omitempty" json:"old_value,omitempty"`
	NewValue    string    `bson:"new_value,omitempty" json:"new_value,omitempty"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}

type WalkoverRecord struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	PlayerID  string    `bson:"player_id" json:"player_id"`
	MatchID   string    `bson:"match_id" json:"match_id"`
	EventID   string    `bson:"event_id" json:"event_id"`
	Reason    string    `bson:"reason" json:"reason"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type MedalRanking struct {
	TournamentID string `bson:"tournament_id" json:"tournament_id"`
	Gold         int    `bson:"gold" json:"gold"`
	Silver       int    `bson:"silver" json:"silver"`
	Bronze       int    `bson:"bronze" json:"bronze"`
	Total        int    `bson:"total" json:"total"`
}
