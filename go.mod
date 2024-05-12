module github.com/rxmeez/chirpy

go 1.22.2

replace github.com/rxmeez/chirpy/internal/database => ./internal/database

require (
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/crypto v0.23.0 // indirect
)
