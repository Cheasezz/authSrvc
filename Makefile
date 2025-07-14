.PHONY: compose-up
compose-up:
	docker-compose -f docker-compose.yml up -d

.PHONY: compose-down
compose-down:
	docker compose down --rmi local --volumes

.PHONY: db-up
db-up:
	docker run --name=test-db -e POSTGRES_PASSWORD=qwerty -p 5432:5432 -d --rm postgres

.PHONY: swag-init
swag-init:
	swag init -g cmd/mainapp/mainapp.go --parseInternal
