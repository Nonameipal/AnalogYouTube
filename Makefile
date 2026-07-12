run:
	go run cmd/main.go

SwagUpdate:
	swag init -g cmd/main.go -o docs --parseInternal

dockUp:
	docker compose up --build -d

dockDown:
	docker compose down

dockLogs:
	docker compose logs -f app

dockRestart:
	docker compose restart app
gofmt:
	gofmt -w cmd internal pkg utils
