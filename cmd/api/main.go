package main

import (
	"log"

	"github.com/Davidmuthee12/socials/internal/env"
	"github.com/Davidmuthee12/socials/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8000"),
	}

	store := store.NewStorage(nil)

	app := &application{
		config: cfg,
		store: store,
	}

	
	mux := app.mount()
	log.Fatal(app.run(mux))
}