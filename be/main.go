package main

import (
	"fmt"
	"log"
	"main/config"
	"main/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()
	config := buildEnv(".env")

	db, err := buildDB(config)
	defer closeDB(db)
	if err != nil {
		log.Printf("Failed to connect to db: %v", err)
		panic(err)
	}
	http.BuildServer(e, db, config)

	if err := (e.Start(fmt.Sprintf(":%s", config.PORT))); err != nil {
		e.Logger.Fatal(err)
	}
}

func buildEnv(env string) *config.Config {
	cfg, err := config.New(env)
	if err != nil {
		panic(err)
	}
	return cfg
}

func buildDB(cfg *config.Config) (*gorm.DB, error) {

	maxIdleConns := 10
	maxOpenConns := 20
	maxLifetime := 15 * time.Minute

	dbCfg := cfg.DB
	sqlCfg := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Database,
	)

	log.Printf("Connecting to db: %s", sqlCfg)

	db, err := gorm.Open(postgres.Open(sqlCfg), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(maxLifetime)
	_, err = sqlDB.Exec(`set search_path='public'`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func closeDB(db *gorm.DB) {
	if db == nil {
		return
	}
	conn, err := db.DB()
	if err != nil {
		log.Printf("Failed to get db connection: %v", err)
		return
	}

	if conn == nil {
		return
	}

	err = conn.Close()
	if err != nil {
		log.Printf("Failed to close db connection: %v", err)
		return
	}
}
