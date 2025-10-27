GeoIP Lookup Service
A self-hosted IP geolocation service built with Go, Redis, and MaxMind databases. This service provides RESTful JSON APIs for looking up geolocation data for single or multiple IP addresses.

Features
🌍 IP geolocation using MaxMind databases

⚡ Redis caching for improved performance

🔄 Automatic database updates from AWS S3

🐳 Environment-based configuration

📊 Structured logging with Zerolog

🔍 Batch IP lookup support

🩺 Health check endpoint

Prerequisites
Go 1.21+

Redis server

AWS account with S3 access

MaxMind account with database access

Quick Start
1. Clone and Build
bash
git clone <your-repo>
cd geoip-service
go mod tidy
go build -o geoip-service
2. Environment Configuration
Create a .env file:

bash
# Server Configuration
SERVER_PORT=8080

# Redis Configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=your_redis_password
REDIS_DB=0

# AWS Configuration
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_aws_access_key_id
AWS_SECRET_ACCESS_KEY=your_aws_secret_access_key
S3_BUCKET=your-geoip-bucket
S3_KEY=databases/GeoLite2-City.mmdb

# Local storage
LOCAL_DB_PATH=./geoip.mmdb
Or set environment variables directly:

bash
export SERVER_PORT=8080
export REDIS_ADDR=localhost:6379
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export S3_BUCKET=your-bucket-name
export S3_KEY=path/to/database.mmdb
3. Prepare S3 Bucket
Upload your MaxMind database to S3:

bash
aws s3 cp GeoLite2-City.mmdb s3://your-bucket-name/databases/GeoLite2-City.mmdb
4. Run the Service
bash
./geoip-service
Or run directly with Go:

bash
go run main.go
API Endpoints
Health Check
Endpoint: GET /health

Check if the service is running properly.

Response:

json
{
  "status": "healthy",
  "timestamp": "2023-10-05T12:00:00Z",
  "service": "geoip-api"
}
Single IP Lookup
Endpoint: GET /lookup/{ip}

Get geolocation data for a single IP address.

Example Request:

bash
curl http://localhost:8080/lookup/8.8.8.8
Response:

json
{
  "ip": "8.8.8.8",
  "country": "United States",
  "country_iso": "US",
  "city": "Mountain View",
  "timezone": "America/Los_Angeles",
  "latitude": 37.386,
  "longitude": -122.0838,
  "accuracy_radius": 1000
}
Batch IP Lookup
Endpoint: POST /batch

Get geolocation data for multiple IP addresses in a single request.

Example Request:

bash
curl -X POST http://localhost:8080/batch \
  -H "Content-Type: application/json" \
  -d '{
    "ips": ["8.8.8.8", "1.1.1.1", "invalid-ip"]
  }'
Response:

json
{
  "8.8.8.8": {
    "ip": "8.8.8.8",
    "country": "United States",
    "country_iso": "US",
    "city": "Mountain View",
    "timezone": "America/Los_Angeles",
    "latitude": 37.386,
    "longitude": -122.0838,
    "accuracy_radius": 1000
  },
  "1.1.1.1": {
    "ip": "1.1.1.1",
    "country": "Australia",
    "country_iso": "AU",
    "city": "Sydney",
    "timezone": "Australia/Sydney",
    "latitude": -33.494,
    "longitude": 143.2104,
    "accuracy_radius": 1000
  },
  "invalid-ip": {
    "ip": "invalid-ip",
    "error": "invalid IP address: invalid-ip"
  }
}
Response Fields
Field	Type	Description
ip	string	The IP address that was queried
country	string	Country name in English
country_iso	string	ISO country code (2 characters)
city	string	City name in English
timezone	string	Time zone (e.g., "America/New_York")
latitude	float64	Latitude coordinate
longitude	float64	Longitude coordinate
accuracy_radius	uint16	Estimated accuracy radius in meters
error	string	Error message (only present on failure)
Configuration Details
Environment Variables
Variable	Required	Default	Description
SERVER_PORT	No	8080	Port the HTTP server listens on
REDIS_ADDR	No	localhost:6379	Redis server address
REDIS_PASSWORD	No	``	Redis password (if any)
REDIS_DB	No	0	Redis database number
AWS_REGION	Yes	-	AWS region for S3 bucket
AWS_ACCESS_KEY_ID	Yes	-	AWS access key ID
AWS_SECRET_ACCESS_KEY	Yes	-	AWS secret access key
S3_BUCKET	Yes	-	S3 bucket name containing the database
S3_KEY	Yes	-	S3 object key for the .mmdb file
LOCAL_DB_PATH	No	./geoip.mmdb	Local path to cache the database
Database Updates
The service downloads the MaxMind database from S3 on startup. To update the database:

Upload a new .mmdb file to your S3 bucket

Restart the service to load the new database

For automatic updates, consider setting up a cron job to:

Download updated databases from MaxMind

Upload to S3

Trigger service restart (or implement hot-reloading)

Docker Deployment
Create a Dockerfile:

dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o geoip-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/geoip-service .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./geoip-service"]
Build and run:

bash
docker build -t geoip-service .
docker run -p 8080:8080 --env-file .env geoip-service
Monitoring and Logging
The service provides structured JSON logs with the following information:

Request method and URL

Response status code

Request duration

Client IP address

User agent

Example log entry:

json
{
  "level": "info",
  "method": "GET",
  "url": "/lookup/8.8.8.8",
  "remote_addr": "192.168.1.100:54321",
  "status": 200,
  "user_agent": "curl/7.68.0",
  "duration": 2.145833,
  "time": "2023-10-05T12:00:00Z"
}
Error Handling
The service returns appropriate HTTP status codes:

200 - Successful lookup

400 - Bad request (invalid IP, malformed JSON)

404 - IP data not found in database

500 - Internal server error

Performance Considerations
Redis cache with 24-hour TTL reduces database lookups

Batch endpoint processes up to 100 IPs per request

Connection pooling for Redis and S3

Efficient memory usage with streaming S3 downloads

License
This project is licensed under the MIT License. Note that MaxMind database usage is subject to MaxMind's license agreement.

Support
For issues and questions:

Check the service logs for error details

Verify Redis connectivity

Confirm AWS credentials and S3 permissions

Ensure the MaxMind database file is valid and accessible
