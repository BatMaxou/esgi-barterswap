compose := docker compose

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'
.PHONY: help

install: ## Construit les images et installe les dépendances Go
	$(compose) build
	$(compose) run --rm go go mod tidy
.PHONY: install

up: ## Démarre l'API et la base (en arrière-plan)
	$(compose) up -d
.PHONY: up

down: ## Arrête et supprime les conteneurs
	$(compose) down
.PHONY: down

build: ## (Re)construit les images Docker
	$(compose) build
.PHONY: build

run: ## Démarre l'API au premier plan (logs visibles)
	$(compose) up
.PHONY: run

logs: ## Affiche les logs du service go
	$(compose) logs -f go
.PHONY: logs

sh: ## Ouvre un shell dans le conteneur go
	$(compose) exec go bash
.PHONY: sh

test: ## Lance les tests
	$(compose) run --rm go go test -v ./...
.PHONY: test

cover: ## Lance les tests avec la couverture
	$(compose) run --rm go go test -cover ./...
.PHONY: cover

fmt: ## Formate le code
	$(compose) run --rm go gofmt -w .
.PHONY: fmt

vet: ## Analyse statique du code
	$(compose) run --rm go go vet ./...
.PHONY: vet

tidy: ## Met à jour go.mod / go.sum
	$(compose) run --rm go go mod tidy
.PHONY: tidy
