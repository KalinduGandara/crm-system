package bootstrap

import "github.com/KalinduGandara/crm-system/db"

type Application struct {
	Env *Env
	DB  db.Client
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.DB = NewMongoDatabase(app.Env)
	return *app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.DB)
}
