package main

import (
	"shb/internal/application"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	a := application.NewApplication()
	a.Start()
}
