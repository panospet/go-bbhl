.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags='-w -s -extldflags "-static"' -o ./bin/highlights cmd/highlights/main.go

.PHONY: container
container: ## create docker container
	docker build -t p4nospet/basketball-highlights .


.PHONY: run-container-euroleague
run-container-euroleague:
	docker run --rm --env-file .env p4nospet/basketball-highlights highlights -euroleague

.PHONY: run-container-nba
run-container-nba:
	docker run --rm --env-file .env p4nospet/basketball-highlights highlights -nba