package pkg

import (
	"Auth/config"
	"Auth/internal/model"
	"context"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserGormRepository struct {
	db *gorm.DB
}

func InitializeDatabase(cfg *config.Config) *UserGormRepository {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DbName)

	fmt.Println(dsn)
	fmt.Println("Connecting to database with DSN:", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("could not connect to the database: %v", err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("failed migration: %v", err)
	}

	return &UserGormRepository{db: db}
}

func (r *UserGormRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserGormRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserGormRepository) Delete(username string) error {
	result := r.db.Where("username = ?", username).Unscoped().Delete(&model.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
