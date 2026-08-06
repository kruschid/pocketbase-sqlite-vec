package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	"github.com/mattn/go-sqlite3"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
)

const (
	driverName    = "pocketbase_sqlite_vec"
	dataDir       = "/pb/pb_data"
	hooksDir      = "/pb/pb_hooks"
	migrationsDir = "/pb/pb_migrations"
	publicDir     = "/pb/pb_public"
)

func init() {
	vec.Auto()

	sql.Register(driverName, &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			_, err := conn.Exec(`
				PRAGMA busy_timeout = 10000;
				PRAGMA journal_mode = WAL;
				PRAGMA journal_size_limit = 200000000;
				PRAGMA synchronous = NORMAL;
				PRAGMA foreign_keys = ON;
				PRAGMA temp_store = MEMORY;
				PRAGMA cache_size = -32000;
			`, nil)
			return err
		},
	})

	dbx.BuilderFuncMap[driverName] = dbx.BuilderFuncMap["sqlite3"]
}

func connectSQLiteVec(path string) (*dbx.DB, error) {
	return dbx.Open(driverName, path)
}

func main() {
	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: dataDir,
		DBConnect:      connectSQLiteVec,
	})

	jsvm.MustRegister(app, jsvm.Config{
		HooksDir:      hooksDir,
		MigrationsDir: migrationsDir,
		HooksWatch:    false,
		HooksPoolSize: 15,
	})

	app.OnServe().BindFunc(func(event *core.ServeEvent) error {
		if !event.Router.HasRoute(http.MethodGet, "/{path...}") {
			event.Router.GET("/{path...}", apis.Static(os.DirFS(publicDir), true))
		}
		return event.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
