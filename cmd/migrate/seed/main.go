package main

import (
	"log"

	"github.com/Davidmuthee12/socials/internal/db"
	"github.com/Davidmuthee12/socials/internal/env"
	"github.com/Davidmuthee12/socials/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://user:adminpassword@localhost/social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	store := store.NewStorage(conn)

	db.Seed(store)
}