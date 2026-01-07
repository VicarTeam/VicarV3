package db

import (
	"context"
	"fmt"
	"os"
	"time"
	"vicar-backend/db/entities"
	"vicar-backend/db/seeder"
	"vicar-backend/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize() error {
	log.Info(log.Database, "🔌", "Connecting to the database...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	log.Success(log.Database, "✅", "Connected to the database.")

	DB = db

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		return err
	}

	models := []any{
		entities.User{},
		entities.RefreshToken{},
		entities.Fido2Login{},
		entities.BaseItem{},
		entities.CharacterDirectory{},
		entities.V5Book{},
		entities.V5BloodRitual{},
		entities.V5OblivionCeremony{},
		entities.V5Clan{},
		entities.V5Discipline{},
		entities.V5DisciplineAbility{},
		entities.V5PredatorType{},
		entities.V5PredatorAction{},
		entities.V5TraitPack{},
		entities.V5Trait{},
		entities.V5TraitPackTrait{},
		entities.V5BloodPotencyData{},
		entities.V5Character{},
		entities.V5CharacterCategory{},
		entities.V5CharacterAttribute{},
		entities.V5CharacterSkill{},
		entities.V5RequiredPointSpread{},
		entities.V5CharacterDisciplineSelection{},
		entities.V5CharacterDisciplineAbility{},
		entities.V5CharacterTraitPackUsage{},
		entities.V5CharacterTrait{},
		entities.V5CharacterFlawTrait{},
		entities.V5CharacterInventory{},
		entities.V5ItemStack{},
		entities.V5Item{},
		entities.V5LevelChange{},
	}

	log.Info(log.Database, "🔄", "Running database migrations...")

	if err := db.AutoMigrate(models...); err != nil {
		return err
	}

	log.Success(log.Database, "✅", "Database migrations completed.")

	log.Info(log.Database, "🌱", "Seeding database with old Vicar data...")
	v5seed := seeder.NewV5Seeder("assets/VicarData")
	if err := v5seed.Seed(context.Background(), db); err != nil {
		return err
	}

	log.Success(log.Database, "✅", "Database seeding completed.")

	return nil
}
