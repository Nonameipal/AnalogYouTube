module github.com/Nonameipal/AnalogYouTube

go 1.25.0

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.8.1
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.12.3
	github.com/redis/go-redis/v9 v9.19.0
	github.com/rs/zerolog v1.34.0
	golang.org/x/crypto v0.50.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
)

replace github.com/Nonameipal/AnalogYouTube => ./
