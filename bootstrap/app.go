package bootstrap

import (
	"be/internal"

	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type State string

const (
	FREEZE State = "freeze"
	ACTIVE State = "active"
	DONE State = "done"
	SETUP State = "setup"
)

type App struct {
	Config     *Config
	State	State
	DB         *gorm.DB
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
	app.DB = NewConnection(app.Config)

	// create token maker
	tokenMaker, err := internal.CreatePasetoMaker(app.Config.JWTSecret)
	if err != nil {
		app.Logger.Error().Err(err).Msg("cannot create token maker")
		os.Exit(1)
	}
	app.TokenMaker = tokenMaker

	app.State = State(app.Config.AppInitialState)
	viper.Set("APP_INITIAL_STATE", app.State)

	return &app
}
