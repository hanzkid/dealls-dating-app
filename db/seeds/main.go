package main

import (
	"log"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/danvergara/seeder"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Seed struct {
	db *sqlx.DB
}

func NewSeed(db *sqlx.DB) Seed {
	return Seed{
		db: db,
	}
}

const (
	connectionString = "postgres://postgres:postgres@postgres:7002/datingapp?sslmode=disable&connect_timeout=5"
)

func main() {
	// Code here
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	s := NewSeed(db)

	if err := seeder.Execute(s); err != nil {
		log.Fatalf("error seeding the db %s\n", err)
	}
}

func (s Seed) UserSeed() {
	count := 50
	password := "$2a$10$RBWbaQe1Ut33bVRw6BTDyOX2oCgHZY3LkjrXO9JZ5qosOINNzgKi2" // 12345678
	for i := 0; i < count; i++ {
		var userID int
		err := s.db.QueryRow(`INSERT INTO public.users(email, name, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`, faker.Email(), faker.Name(), password, time.Now(), time.Now()).Scan(&userID)
		if err != nil {
			log.Fatalf("error seeding user: %v", err)
		}
		_, err = s.db.Exec(`INSERT INTO public.profiles(user_id, description, picture, created_at, updated_at) VALUES($1, $2, $3, $4, $5)`, userID, faker.Sentence(), "https://picsum.photos/id/1/200/300", time.Now(), time.Now())
	}
}
