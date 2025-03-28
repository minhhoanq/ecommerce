up-dev:
	@echo "starting containers..."
	docker compose -f docker-compose.dev.yaml up --build --remove-orphans

.PHONY: up-dev