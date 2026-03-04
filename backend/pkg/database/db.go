package database

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	interval2  = 2
	interval10 = 10
	interval20 = 20
	interval30 = 30
)

func Connect(mongoURI string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetMaxPoolSize(interval20).
		SetMinPoolSize(interval2).
		SetMaxConnIdleTime(interval30 * time.Second).
		SetServerSelectionTimeout(interval10 * time.Second).
		SetHeartbeatInterval(interval10 * time.Second).
		SetRetryWrites(true).
		SetRetryReads(true).
		SetReadConcern(readconcern.Majority())

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to mongodb")
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Error().Err(err).Msg("failed to ping mongodb")
		return nil, err
	}
	log.Info().Msg("connected to mong db")
	return client, nil
}
