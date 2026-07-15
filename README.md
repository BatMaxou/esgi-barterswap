# BarterSwap — Skill exchange API

Time bank: each hour of service provided grants one hour of service received.
Exchanges are settled in **time-credits**.

## Installation

```bash
git clone git@github.com:BatMaxou/esgi-barterswap.git
cd esgi-barterswap
make install   # build images + go mod tidy
make up        # start the API (http://localhost:8000) and the database
```

The project is driven through the `Makefile` (`make help` for the list of
targets): `make test`, `make cover`, `make fmt`, `make vet`, `make logs`,
`make down`.

## Endpoints

| Method | Path | Auth | Description |
|--------|------|:----:|-------------|
| `GET`    | `/`                     |   | Service status |
| `POST`   | `/api/users`            |   | Create an account (10 welcome credits granted automatically) |
| `GET`    | `/api/users/{id}`       |   | Public profile of a user (balance computed from the ledger) |
| `PUT`    | `/api/users/{id}`       | ✔ | Update your own profile |
| `GET`    | `/api/users/{id}/skills`|   | Skills of a user |
| `PUT`    | `/api/users/{id}/skills`| ✔ | Set your skills (overwritten on each call) |
| `GET`    | `/api/services`         |   | List of ads (filters `category`, `city`, `search`) |
| `POST`   | `/api/services`         | ✔ | Publish a service ad |
| `GET`    | `/api/services/{id}`    |   | Ad details |
| `PUT`    | `/api/services/{id}`    | ✔ | Update your own ad |
| `DELETE` | `/api/services/{id}`    | ✔ | Delete your own ad |

Authentication is done through the `X-User-ID: {id}` header.
Service filtering and search are performed **server-side** (query params).

## Usage examples

Create a user (success → `201`):

```bash
curl -i -X POST http://localhost:8000/api/users \
  -H 'Content-Type: application/json' \
  -d '{"pseudo":"Thierry","bio":"my bio","city":"Paris"}'
```

```json
{
  "id": 1,
  "pseudo": "Thierry",
  "bio": "my bio",
  "city": "Paris",
  "credit_balance": 10,
  "created_at": "2026-06-25T19:55:22Z"
}
```

Empty pseudo → `400`:

```bash
curl -i -X POST http://localhost:8000/api/users \
  -H 'Content-Type: application/json' \
  -d '{"pseudo":""}'
# {"error":"pseudo is required"}
```

Get a user profile (success → `200`):

```bash
curl -i http://localhost:8000/api/users/1
# 200 -> {"id":1,"pseudo":"Thierry",...,"credit_balance":10,...}
# 404 -> {"error":"user not found"}        (unknown id)
# 400 -> {"error":"invalid identifier"}    (non-numeric id)
```

Publish a service ad (auth required → `201`):

```bash
curl -i -X POST http://localhost:8000/api/services \
  -H 'Content-Type: application/json' \
  -H 'X-User-ID: 1' \
  -d '{"title":"Cours de Go","description":"Initiation en 1h","category":"Informatique","duration_minutes":60,"credits":2,"city":"Paris"}'
# 201 -> {"id":1,"provider_id":1,"title":"Cours de Go","category":"Informatique","active":true,...}
# 400 -> {"error":"invalid category (not in the category list)"}
# 401 -> {"error":"authentication required (X-User-ID header)"}
```

Search ads (server-side filters → `200`):

```bash
curl -i 'http://localhost:8000/api/services?category=Informatique&city=Paris&search=Go'
# 200 -> [{"id":1,"title":"Cours de Go",...}]
```

Delete your own ad (auth required → `204`):

```bash
curl -i -X DELETE http://localhost:8000/api/services/1 -H 'X-User-ID: 1'
# 204 (no content)
# 403 -> {"error":"action not allowed"}   (ad owned by another user)
# 404 -> {"error":"service not found"}
```

## Tests

```bash
make test
make cover
```
