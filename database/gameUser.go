package database

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Team is a string definging whether a team is red or blue
type Team string

const (
	RedTeam  Team = "red"
	BlueTeam Team = "blue"
)

type GameUser struct {
	ID        string `gorm:"unique"`
	Email     string
	Username  string
	Password  string
	CreatedAt time.Time
	GameID    string
	Team      Team
}

// AddUser inserts a new user into the database
func AddUser(ctx context.Context, username, email, password, gameID string, team Team) (GameUser, error) {
	var user GameUser
	user.ID = uuid.New().String()[0:8]
	user.Username = username
	user.Email = email
	user.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	user.Team = team
	user.CreatedAt = time.Now()
	user.GameID = gameID
	if err := pool.WithContext(ctx).Create(&user).Error; err != nil {
		return GameUser{}, err
	}
	return user, nil
}

func AuthUser(ctx context.Context, username, password, gameid string) (GameUser, error) {
	var user GameUser
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	if err := pool.Model(GameUser{}).Where("username = ? AND password = ? AND game_id=?", username, hash, gameid).Find(&user).Error; err != nil {
		return GameUser{}, err
	}
	return user, nil
}
