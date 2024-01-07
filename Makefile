.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	go build -o ./bin/highlights cmd/highlights/main.go

.PHONY: build-sample
build-sample:
	go build -o ./bin/sample cmd/sample/main.go

.PHONY: container
container: ## create docker container
	docker build -t p4nospet/basketball-highlights .

.PHONY: run-euroleague
run-euroleague:
	docker run --rm --env-file .env -v $(PWD)/data:/data p4nospet/basketball-highlights highlights -euroleague

.PHONY: run-nba
run-nba:
	docker run --rm --env-file .env -v $(PWD)/data:/data p4nospet/basketball-highlights highlights -nba

.PHONY: dry-run-euroleague
dry-run-euroleague:
	docker run --rm --env-file .env -v $(PWD)/data:/data p4nospet/basketball-highlights highlights -euroleague -dry

.PHONY: dry-run-nba
dry-run-nba:
	docker run --rm --env-file .env -v $(PWD)/data:/data p4nospet/basketball-highlights highlights -nba -dry

.PHONY: send-sample
run-container-sample:
	docker run --rm --env-file .env p4nospet/basketball-highlights sample