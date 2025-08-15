package main

import (
	"log"

	"github.com/I-Van-Radkov/messenger/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	cfg := config.MustLoad()

}
