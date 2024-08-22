package main

import (
	"database/sql"
	"os"

	"github.com/demola234/defiraise/api"
	db "github.com/demola234/defiraise/db/sqlc"
	"github.com/demola234/defiraise/utils"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// @title          DefiFundr API
// @version         1.0
// @description     Decentralized Crowdfunding Platform for DeFi Projects
// @contact.name   DefiFundr
// @schemes        http
// @contact.url    http://www.swagger.io/support
// @contact.email  kolawoleoluwasegun567@gmail
// @host localhost:8080

func main() {
	configs, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	if configs.Environment == "developement" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}

	conn, err := sql.Open(configs.DBDriver, configs.DBSource)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	store := db.NewStore(conn)

	runGinServer(configs, store)
}

func runGinServer(configs utils.Config, store db.Store) {
	server, err := api.NewServer(configs, store)

	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

	err = server.Start(configs.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot start http server")
	}
}
