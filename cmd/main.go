package main

import (
	"gitlab.com/climate.uz/internal/app"
	"gitlab.com/climate.uz/config"
)

func main() {
	cfg := config.Load()
	app.Run(&cfg)
}
