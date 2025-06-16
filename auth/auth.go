package auth

import (
	"context"
	"errors"
	"time"

	"github.com/appu900/deerDB/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	config         *types.Config
	userCollection *mongo.Collection
	mongoCtx       context.Context
	cancelFunc     context.CancelFunc
}

type UserRegisterResponse struct {
	UserID    string `json:"userID"`
	UserName  string `json:"username"`
	UserEmail string `json:"useremail"`

}

func ToUserResponse(user *types.User) *UserRegisterResponse {
	return &UserRegisterResponse{
		UserID:    user.ID,
		UserName:  user.Username,
		UserEmail: user.Email,
	}
}

func NewAuthService(config *types.Config) (*AuthService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		cancel()
		return nil, err
	}
	userCollection := client.Database("deerdb").Collection("users")
	return &AuthService{
		config:         config,
		userCollection: userCollection,
		mongoCtx:       ctx,
		cancelFunc:     cancel,
	}, nil
}

func (a *AuthService) Close() {
	a.cancelFunc()
}

func (authService *AuthService) RegisterUser(username, email, password string) (*types.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := authService.userCollection.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"username": username},
			{"email": email},
		},
	})
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("user already exists")
	}
	// password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &types.User{
		ID:       primitive.NewObjectID().Hex(),
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Created:  time.Now(),
	}
	_, err = authService.userCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
