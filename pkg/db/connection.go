package db

import (
	"ShowTimes/pkg/config"
	"ShowTimes/pkg/domain"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// we have to AutoMigrate to database to add users and store  information about them

func ConectDatabse(cfg config.Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)
	db, dbErr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{SkipDefaultTransaction: true})

	if err := db.AutoMigrate(&domain.Users{}); err != nil {
		return db, err
	}
	if err := db.AutoMigrate(&domain.Category{}); err != nil {
		return db, err
	}
	if err := db.AutoMigrate(&domain.Product{}); err != nil {
		return db, err
	}
	if err := db.AutoMigrate(&domain.Address{}); err != nil {
		return db, err
	}
	if err := db.AutoMigrate(&domain.Cart{}); err != nil {
		return db, err
	}
	if err := db.AutoMigrate(&domain.Order{}); err != nil {
		return db, err
	}
	if err := db.AutoMigrate(&domain.OrderItem{}); err != nil {
		return db, err
	}

	if err := db.AutoMigrate(&domain.Admin{}); err != nil {
		return db, err
	}

	if err := db.AutoMigrate(&domain.Payment{}); err != nil {
		return db, err
	}
	if err := db.AutoMigrate(&domain.Wallet{}); err != nil {
		return db, err
	}
	if err := db.AutoMigrate(&domain.ProductOffer{}); err != nil {
		return db, err
	}
	if err := db.AutoMigrate(&domain.CategoryOffer{}); err != nil {
		return db, err
	}
	if err := db.AutoMigrate(&domain.Referral{}); err != nil {
		return db, err
	}
	if err := db.AutoMigrate(&domain.Coupon{}); err != nil {
		return db, err
	}

	// CheckAndCreateAdmin(db)
	return db, dbErr
}

func CheckAndCreateAdmin(db *gorm.DB) {
	var count int64
	db.Model(&domain.Users{}).Count(&count)
	if count == 0 {
		password := "admin@123"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return
		}
		admin := domain.Users{
			ID:       1,
			Name:     "Admin ShowTimes",
			Email:    "admin@showtimes.com",
			Password: string(hashedPassword),
			IsAdmin:  true,
			Blocked:  false,
			Phone:    "9746359523",
		}
		db.Create(&admin)
	}
}
