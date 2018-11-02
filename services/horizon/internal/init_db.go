package horizon

import (
	"github.com/stellar/go/services/horizon/internal/db2/core"
	"github.com/stellar/go/services/horizon/internal/db2/history"
	"github.com/stellar/go/support/db"
	"github.com/stellar/go/support/log"
)

func initHorizonDb(app *App) {
	session, err := db.Open("postgres", app.config.DatabaseURL)

	if err != nil {
		log.Panic(err)
	}
	session.DB.SetMaxIdleConns(200)
	session.DB.SetMaxOpenConns(200)

	app.historyQ = &history.Q{session}
}

func initCoreDb(app *App) {
	session, err := db.Open("postgres", app.config.StellarCoreDatabaseURL)

	if err != nil {
		log.Panic(err)
	}

	session.DB.SetMaxIdleConns(200)
	session.DB.SetMaxOpenConns(200)
	app.coreQ = &core.Q{session}
}

func init() {
	appInit.Add("horizon-db", initHorizonDb, "app-context", "log")
	appInit.Add("core-db", initCoreDb, "app-context", "log")
}
