package bootstrap

import (
	"be/db/sqlc"
	"be/internal"

	"os"
	"time"

	"github.com/rs/zerolog"
)

type App struct {
	Config     *Config
	DB         *sqlc.Queries
	Logger     zerolog.Logger
	TokenMaker internal.Maker
}

func NewApp() *App {
	app := App{}

	// initialize logger
	app.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger().With().Caller().Logger()

	// initialize config
	app.Config = NewConfig()

	// connect DB
	conn := NewConnection(app.Config)
	app.DB = sqlc.New(conn)

	// create token maker
	tokenMaker, err := internal.CreatePasetoMaker(app.Config.JWTSecret)
	if err != nil {
		app.Logger.Error().Err(err).Msg("cannot create token maker")
		os.Exit(1)
	}
	app.TokenMaker = tokenMaker

	return &app
}
