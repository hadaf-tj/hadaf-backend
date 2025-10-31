package main

import (
	_ "github.com/joho/godotenv/autoload"
	"shb/internal/application"
)

func main() {
	a := application.NewApplication()
	a.Start()
}
