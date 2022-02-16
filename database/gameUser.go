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
	NoTeam   Team = "no"
)

type GameUser struct {
	ID         string
	Email      string
	Username   string `gorm:"primaryKey"`
	Password   string
	CreatedAt  time.Time
	GameID     string `gorm:"primaryKey"`
	Team       Team
	JoinedGame bool
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
	user.JoinedGame = false
	if err := pool.WithContext(ctx).Create(&user).Error; err != nil {
		return GameUser{}, err
	}
	return user, nil
}

func UpdateUserStart(ctx context.Context, username, gameid string) (GameUser, error) {
	var user GameUser
	user, err := CheckUser(ctx, username, gameid)
	if err != nil {
		return GameUser{}, err
	}
	if err := pool.Model(&user).Update("joined_game", true).Error; err != nil {
		return GameUser{}, err
	}
	return user, nil
}

func UpdateUsersTeam(ctx context.Context, username, gameid string, team Team) (GameUser, error) {
	var user GameUser
	user, err := CheckUser(ctx, username, gameid)
	if err != nil {
		return GameUser{}, err
	}
	if err := pool.Model(&user).Update("team", team).Error; err != nil {
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

func CheckUser(ctx context.Context, username, gameid string) (GameUser, error) {
	var user GameUser
	if err := pool.Model(GameUser{}).Where("username = ? AND game_id=?", username, gameid).Find(&user).Error; err != nil {
		return GameUser{}, err
	}
	return user, nil
}
