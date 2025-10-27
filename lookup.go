package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
)

type GeoIPResponse struct {
	IP         string  `json:"ip"`
	Country    string  `json:"country,omitempty"`
	CountryISO string  `json:"country_iso,omitempty"`
	City       string  `json:"city,omitempty"`
	Timezone   string  `json:"timezone,omitempty"`
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
	Accuracy   uint16  `json:"accuracy_radius,omitempty"`
	Error      string  `json:"error,omitempty"`
}

func (s *GeoIPService) LookupIP(ipStr string) (*GeoIPResponse, error) {
	ctx := context.Background()
	cacheKey := "geoip:" + ipStr

	// 1. Try getting result from Redis cache
	cachedResult, err := s.rdb.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var response GeoIPResponse
		if err := json.Unmarshal(cachedResult, &response); err == nil {
			s.log.Debug().Str("ip", ipStr).Msg("Served from cache")
			return &response, nil
		}
	} else if err != redis.Nil {
		s.log.Warn().Err(err).Str("ip", ipStr).Msg("Redis error")
		// Proceed to database lookup
	}

	// 2. Perform MaxMind database lookup
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	record, err := s.db.City(ip)
	if err != nil {
		return nil, err
	}

	// Build response
	response := &GeoIPResponse{
		IP:         ipStr,
		Country:    getEnglishName(record.Country.Names),
		CountryISO: record.Country.IsoCode,
		City:       getEnglishName(record.City.Names),
		Timezone:   record.Location.TimeZone,
		Accuracy:   record.Location.AccuracyRadius,
	}

	if record.Location.Latitude != 0 || record.Location.Longitude != 0 {
		response.Latitude = record.Location.Latitude
		response.Longitude = record.Location.Longitude
	}

	// 3. Cache the result in Redis
	jsonData, err := json.Marshal(response)
	if err != nil {
		return response, nil // Return response even if caching fails
	}

	if err := s.rdb.Set(ctx, cacheKey, jsonData, 24*time.Hour).Err(); err != nil {
		s.log.Warn().Err(err).Msg("Failed to cache result")
	}

	return response, nil
}

// Helper function to get English name from map
func getEnglishName(names map[string]string) string {
	if name, exists := names["en"]; exists {
		return name
	}
	return ""
}
