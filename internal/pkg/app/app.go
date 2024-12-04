package app

import (
	enpoint "db-worker/internal/endpoint/app"
	dbworker "db-worker/internal/service/db"
	generator "db-worker/internal/service/keygen"
	api "db-worker/internal/transport/rest"
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

	worker := dbworker.New(a.app.GetLogger())
	keygen := generator.New(35)
	api := api.New(worker, keygen)

	a.app.AddUsersGetHandler("/all", api.GetAllUsers)
	a.app.AddUsersPostHandler("/email/", api.GetUser)
	a.app.AddUsersPostHandler("/id/", api.GetUserByID)
	a.app.AddUsersPostHandler("/compare", api.ComparePassword)
	a.app.AddUsersPostHandler("/register", api.RegisterNewUser)
	a.app.AddUsersPostHandler("/check/", api.CheckExists)

	a.app.AddPostHandler("/", api.RegisterOperation) // Регистрация
	a.app.AddGetHandler("/id", api.GetID)            // Получение ID операции
	a.app.AddGetHandler("/", api.GetAllOperations)   // Получение всех операций

	a.app.AddOperationGetHandler("/version/:id", api.GetVersion) // Получение версии
	a.app.AddOperationPostHandler("/:id", api.SetResult)         // Обновление
	a.app.AddOperationGetHandler("/:id", api.GetOperation)       // Получение

	a.app.GetLogger().Info("Starting DB-service")
	a.app.Start()
}
