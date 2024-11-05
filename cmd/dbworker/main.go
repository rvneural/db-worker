package main

import "db-worker/internal/pkg/app"

func main() {
	App := app.New()
	App.Run()
}
