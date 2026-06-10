package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/tt-tournament/backend/internal/middleware"
	"github.com/tt-tournament/backend/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	col *mongo.Collection
}

func NewUserService(db *mongo.Database) *UserService {
	return &UserService{col: db.Collection("users")}
}

func (s *UserService) Register(ctx context.Context, username, password, displayName, role string) (*model.User, error) {
	count, err := s.col.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("用户名已存在")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := model.User{
		ID:          primitive.NewObjectID().Hex(),
		Username:    username,
		Password:    string(hash),
		DisplayName: displayName,
		Role:        role,
		Ranking:     0,
		CreatedAt:   time.Now(),
	}

	_, err = s.col.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Login(ctx context.Context, username, password string) (string, *model.User, error) {
	var user model.User
	err := s.col.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return "", nil, fmt.Errorf("用户名或密码错误")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, fmt.Errorf("用户名或密码错误")
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", nil, err
	}
	return token, &user, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := s.col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateRanking(ctx context.Context, id string, ranking int) error {
	_, err := s.col.UpdateByID(ctx, id, bson.M{"$set": bson.M{"ranking": ranking}})
	return err
}

func (s *UserService) ListByRole(ctx context.Context, role string) ([]model.User, error) {
	filter := bson.M{}
	if role != "" {
		filter["role"] = role
	}
	cursor, err := s.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var users []model.User
	cursor.All(ctx, &users)
	return users, nil
}

func GenerateFileHash(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
