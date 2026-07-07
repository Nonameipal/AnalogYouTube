run:
	go run cmd/main.go

SwagUpdate:
	swag init -g cmd/main.go -o docs --parseInternal