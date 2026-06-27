# BarterSwap — API d'échange de compétences

Banque de temps : chaque heure de service rendue donne droit à une heure de
service reçue. Les échanges sont réglés en **crédits-temps**.

## Installation

```bash
git clone git@github.com:BatMaxou/esgi-barterswap.git
cd esgi-barterswap
make install   # build des images + go mod tidy
make up        # démarre l'API (http://localhost:8000) et la base
```

Le projet se pilote via le `Makefile` (`make help` pour la liste des cibles) :
`make test`, `make cover`, `make fmt`, `make vet`, `make logs`, `make down`.

## Endpoints

| Méthode | Path | Description |
|---------|------|-------------|
| `GET`  | `/`               | État du service |
| `POST` | `/api/users`      | Créer un compte (10 crédits de bienvenue attribués automatiquement) |
| `GET`  | `/api/users/{id}` | Profil public d'un utilisateur (solde calculé depuis le journal) |

## Exemples d'utilisation

Créer un utilisateur (succès → `201`) :

```bash
curl -i -X POST http://localhost:8000/api/users \
  -H 'Content-Type: application/json' \
  -d '{"pseudo":"Thierry","bio":"ma bio","ville":"Paris"}'
```

```json
{
  "id": 1,
  "pseudo": "Thierry",
  "bio": "ma bio",
  "ville": "Paris",
  "credit_balance": 10,
  "created_at": "2026-06-25T19:55:22Z"
}
```

Pseudo vide → `400` :

```bash
curl -i -X POST http://localhost:8000/api/users \
  -H 'Content-Type: application/json' \
  -d '{"pseudo":""}'
# {"error":"le pseudo est obligatoire"}
```

Récupérer le profil d'un utilisateur (succès → `200`) :

```bash
curl -i http://localhost:8000/api/users/1
# 200 -> {"id":1,"pseudo":"Thierry",...,"credit_balance":10,...}
# 404 -> {"error":"utilisateur introuvable"}   (id inexistant)
# 400 -> {"error":"identifiant invalide"}        (id non numérique)
```

## Tests

```bash
make test    # go test -v ./...
make cover   # go test -cover ./...
```
