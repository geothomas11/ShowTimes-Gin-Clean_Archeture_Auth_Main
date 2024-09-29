package db

import (
	"ShowTimes/pkg/config"
	"ShowTimes/pkg/domain"
	"fmt"

	// "golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

	// if err := db.AutoMigrate(&domain.Admin{}); err != nil {
	// 	return db, err
	// }
	// CheckAndCreateAdmin(db)
	return db, dbErr
}

// func CheckAndCreateAdmin(db *gorm.DB) {
// 	var count int64
// 	db.Model(&domain.Users{}).Count(&count)
// 	if count == 0 {
// 		password := "admin@123"
// 		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 		fmt.Println(hashedPassword)
// 		if err != nil {
// 			return
// 		}
// 		admin := domain.Users{
// 			ID:       1,
// 			Name:     "Showtimes",
// 			Email:    "showtimes@showtimes.com",
// 			Password: string(hashedPassword),
// 			Phone:    "1234567890",
// 			Blocked:  false,
// 			IsAdmin:  true,
// 		}
// 		db.Create(&admin)
// 	}

// }
