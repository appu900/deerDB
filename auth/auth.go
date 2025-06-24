package auth

import (
	"context"
	"errors"
	"time"

	"github.com/appu900/deerDB/types"
	"github.com/golang-jwt/jwt/v5"
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

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
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

func (authService *AuthService) Login(email, password string) (*types.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user types.User
	err := authService.userCollection.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&user)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("Invalid Password")
	}
	return &user, nil
}

func (authService *AuthService) GenerateToken(user *types.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(authService.config.JWTSecret))
}

func (authService *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(authService.config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (authService *AuthService) GetUserByID(userID string) (*types.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user types.User
	err := authService.userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

