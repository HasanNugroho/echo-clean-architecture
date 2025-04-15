package configs

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (cfg *DatabaseConfig) InitMongo(logger *zerolog.Logger) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	monitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			Logger.Info().Msgf("[MONGO] CMD: %s | %s => %v\n", evt.CommandName, evt.Command, evt.RequestID)
		},
		Succeeded: func(_ context.Context, evt *event.CommandSucceededEvent) {
			Logger.Info().Msgf("[MONGO] CMD SUCCESS: %s | Duration: %v\n", evt.CommandName, evt.Duration)
		},
		Failed: func(_ context.Context, evt *event.CommandFailedEvent) {
			Logger.Info().Msgf("[MONGO] CMD FAIL: %s | Error: %v\n", evt.CommandName, evt.Failure)
		},
	}

	clientOpts := options.Client().
		ApplyURI(cfg.URI).
		SetMonitor(monitor)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, err
	}

	// Ping to confirm connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	logger.Info().Msg("✅ Database MongoDB connected successfully!")
	return client.Database(cfg.Database), nil
}
