package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/oschwald/geoip2-golang"
	"github.com/redis/go-redis/v9"
)

type GeoIPService struct {
	db     *geoip2.Reader
	rdb    *redis.Client
	log    *zerolog.Logger
	config *Config
}

func NewGeoIPService(ctx context.Context, cfg *Config) (*GeoIPService, error) {
	logger := zerolog.Ctx(ctx)

	// Download the database from S3
	if err := downloadFromS3(ctx, cfg); err != nil {
		return nil, fmt.Errorf("failed to download database from S3: %w", err)
	}

	// Open the MaxMind database
	db, err := geoip2.Open(cfg.LocalDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open MaxMind database: %w", err)
	}

	// Connect to Redis
	redisOpt := &redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}

	rdb := redis.NewClient(redisOpt)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info().Msg("GeoIP service initialized successfully")

	return &GeoIPService{
		db:     db,
		rdb:    rdb,
		log:    logger,
		config: cfg,
	}, nil
}

func (s *GeoIPService) Close() {
	if s.db != nil {
		s.db.Close()
	}
	if s.rdb != nil {
		s.rdb.Close()
	}
}
