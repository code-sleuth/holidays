package main

import (
	"fmt"
	"holidays/pkg/database"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	mongoDbURI, databaseName, port string
)

func init() {
	viper.AutomaticEnv()
	getEnvVars()
}

func getEnvVars() {
	mongoDbURI = viper.GetString("MONGO_DB_URL")
	if strings.EqualFold(mongoDbURI, "") {
		mongoDbURI = "mongodb://localhost:27017"
	}
	databaseName = viper.GetString("DATABASE_NAME")
	if strings.EqualFold(databaseName, "") {
		databaseName = "holidays"
	}
	port = viper.GetString("PORT")
	if strings.EqualFold(port, "") {
		port = "8080"
	}
}

func main() {
	db, err := database.Connect(mongoDbURI)
	if err != nil {
		log.Info().Msgf("mongoDbURI: %+v", mongoDbURI)
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	e := buildWebService(db, databaseName)

	log.Info().Str("port", port).Msg("starting server")
	if err := e.Start(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
}
