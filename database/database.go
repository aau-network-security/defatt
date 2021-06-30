package database

import (
	"context"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	pool *gorm.DB
)

func New(ctx context.Context, fileName string) error {
	db, err := gorm.Open(sqlite.Open(fileName), &gorm.Config{})
	if err != nil {
		return err
	}
	pool = db
	if err := initalMigration(ctx); err != nil {
		return err
	}
	return nil
}

func Close() error {
	sqlDB, err := pool.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}
	pool = nil
	return nil
}

func initalMigration(ctx context.Context) error {
	if err := pool.WithContext(ctx).AutoMigrate(&GameUser{}); err != nil {
		return err
	}
	return nil
}
