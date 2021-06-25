package database

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Team string

const (
	RedTeam  Team = "red"
	BlueTeam Team = "blue"
)

type GameUser struct {
	ID       string
	Email    string
	Username string
	Password string
	// Metadata  map[string]string
	CreatedAt time.Time
	GameID    string
	Team      Team
}

func AddUser(ctx context.Context, username, email, password string, team Team) (GameUser, error) {
	var user GameUser
	user.ID = uuid.New().String()[0:8]
	user.Username = username
	user.Email = email
	user.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	user.Team = team
	user.CreatedAt = time.Now()
	if err := pool.WithContext(ctx).Create(&user).Error; err != nil {
		return GameUser{}, err
	}
	return user, nil
}

func AuthUser(ctx context.Context, username, password, gameid string) (GameUser, error) {
	var user GameUser
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	if err := pool.Model(GameUser{}).Where("name = ? AND password = ?", username, hash).Find(&user).Error; err != nil {
		return GameUser{}, err
	}
	return user, nil
}
