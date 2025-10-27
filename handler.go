package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (s *GeoIPService) SingleIPHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ipStr := vars["ip"]

	if ipStr == "" {
		http.Error(w, `{"error": "Missing IP address"}`, http.StatusBadRequest)
		return
	}

	result, err := s.LookupIP(ipStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(GeoIPResponse{
			IP:    ipStr,
			Error: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *GeoIPService) BatchIPHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		IPs []string `json:"ips"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if len(request.IPs) == 0 {
		http.Error(w, `{"error": "No IP addresses provided"}`, http.StatusBadRequest)
		return
	}

	// Limit batch size
	if len(request.IPs) > 100 {
		http.Error(w, `{"error": "Batch size too large (max 100 IPs)"}`, http.StatusBadRequest)
		return
	}

	results := make(map[string]GeoIPResponse)
	for _, ip := range request.IPs {
		data, err := s.LookupIP(ip)
		if err != nil {
			results[ip] = GeoIPResponse{
				IP:    ip,
				Error: err.Error(),
			}
		} else {
			results[ip] = *data
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// Health check handler
func (s *GeoIPService) HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "geoip-api",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
