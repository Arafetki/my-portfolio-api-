package main

import (
	"log/slog"
	"os"
	"runtime/debug"
	"sync"

	database "github.com/Arafetki/my-portfolio-api/internal/db"
	"github.com/Arafetki/my-portfolio-api/internal/env"
	"github.com/lmittmann/tint"
)

// @title API Documentation
// @description Personal Portfolio Rest API
// @version 1.0.0
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /v1

type config struct {
	httpPort int
	env      string
	db       struct {
		dsn         string
		automigrate bool
	}
}

type application struct {
	cfg    config
	logger *slog.Logger
	wg     sync.WaitGroup
}

const version = "1.0.0"

func main() {

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {

	var cfg config
	cfg.httpPort = env.GetInt("APP_PORT", 8080)
	cfg.env = env.GetString("APP_ENV", "development")
	cfg.db.automigrate = env.GetBool("DB_AUTOMIGRATE", true)

	// // Establish Vault Connection
	// v, err := vault.NewVault("secret")
	// if err != nil {
	// 	return err
	// }
	// // Initialize Secret Store
	// secretStore := secrets.NewStore(v)

	cfg.db.dsn = env.GetString("DB_DSN", "")
	db, err := database.New(cfg.db.dsn, cfg.db.automigrate)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("db connection has been established sucessfully!")

	app := &application{
		cfg:    cfg,
		logger: logger,
	}

	return app.serveHTTP()
}
