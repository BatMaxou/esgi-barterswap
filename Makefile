compose := docker compose

export USER_ID := $(shell id -u)
export GROUP_ID := $(shell id -g)

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

reset: ## Détruit la base (volume) et la recrée vide
	$(compose) down -v
	${MAKE} up
.PHONY: reset

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
	$(compose) run --rm go go test -v ./src
.PHONY: test

cover: ## Lance les tests avec la couverture
	$(compose) run --rm go go test -cover ./src
.PHONY: cover

fmt: ## Formate le code
	$(compose) run --rm go gofmt -w ./src
.PHONY: fmt

fmt-check: ## Vérifie le formatage du code
	$(compose) run --rm go gofmt -l ./src
.PHONY: fmt-check

vet: ## Analyse statique du code
	$(compose) run --rm go go vet ./src
.PHONY: vet

tidy: ## Met à jour go.mod / go.sum
	$(compose) run --rm go go mod tidy
.PHONY: tidy

exec-db: ## Ouvre une console dans le conteneur DB
	$(compose) exec db mariadb -uroot -proot barterswap
.PHONY: exec-db
