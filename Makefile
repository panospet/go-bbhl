.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags='-w -s -extldflags "-static"' -o ./bin/highlights cmd/highlights/main.go

.PHONY: container
container: ## create docker container
	docker build -t p4nospet/basketball-highlights .


.PHONE: run-container-euroleague
run-container-euroleague:
	docker run -it --rm --env-file .env p4nospet/basketball-highlights highlights -euroleague

.PHONE: run-container-nba
run-container-nba:
	docker run -it --rm --env-file .env p4nospet/basketball-highlights highlights -nba