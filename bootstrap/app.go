package bootstrap

import "github.com/KalinduGandara/crm-system/db"

type Application struct {
	Env *Env
	DB  db.Client
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.DB = NewMySQLDatabase(app.Env)
	return *app
}

func (app *Application) CloseDBConnection() {
	CloseMySQLDBConnection(app.DB)
}
