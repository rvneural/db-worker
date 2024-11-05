package app

import (
	enpoint "db-worker/internal/endpoint/app"
)

type App struct {
	app *enpoint.Endpoint
}

func New() *App {
	return &App{
		app: enpoint.New(),
	}
}

func (a *App) Run() {
	a.app.AddPostHandler("/", nil) // Регистрация
	a.app.AddGetHandler("/", nil)  // Получение всех операций

	a.app.AddOperationPostHandler("/:id", nil)        // Обновление
	a.app.AddOperationGetHandler("/version/:id", nil) // Получение версии
	a.app.AddOperationGetHandler("/:id", nil)         // Получение

	a.Run()
}
