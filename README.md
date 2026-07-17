# BarterSwap — Skill exchange API

Time bank: each hour of service provided grants one hour of service received.
Exchanges are settled in **time-credits**.

## Installation

```bash
git clone git@github.com:BatMaxou/esgi-barterswap.git
cd esgi-barterswap
make install
make up
```

The project is driven through the `Makefile` (`make help` for the list of
targets): `make test`, `make cover`, `make fmt`, `make vet`, `make logs`,
`make down`.

## Endpoints

| Method | Path | Auth | Description |
|--------|------|:----:|-------------|
| `GET`    | `/`                            |   | Service status |
| `POST`   | `/api/users`                   |   | Create an account (10 welcome credits granted automatically) |
| `GET`    | `/api/users/{id}`              |   | Public profile of a user (balance computed from the ledger) |
| `PUT`    | `/api/users/{id}`              | ✔ | Update your own profile |
| `GET`    | `/api/users/{id}/skills`       |   | Skills of a user |
| `PUT`    | `/api/users/{id}/skills`       | ✔ | Set your skills (overwritten on each call) |
| `GET`    | `/api/users/{id}/reviews`      |   | Reviews received by a user |
| `GET`    | `/api/users/{id}/stats`        |   | Dashboard of a user |
| `GET`    | `/api/services`                |   | List of ads (filters `category`, `city`, `search`) |
| `POST`   | `/api/services`                | ✔ | Publish a service ad |
| `GET`    | `/api/services/{id}`           |   | Ad details |
| `PUT`    | `/api/services/{id}`           | ✔ | Update your own ad |
| `DELETE` | `/api/services/{id}`           | ✔ | Delete your own ad |
| `GET`    | `/api/services/{id}/reviews`   |   | Reviews on a service |
| `POST`   | `/api/exchanges`               | ✔ | Request an exchange on a service |
| `GET`    | `/api/exchanges`               | ✔ | Your exchanges, requested and received (filter `status`) |
| `GET`    | `/api/exchanges/{id}`          | ✔ | Exchange details (participants only) |
| `PUT`    | `/api/exchanges/{id}/accept`   | ✔ | Accept a request (owner) — credits are held |
| `PUT`    | `/api/exchanges/{id}/reject`   | ✔ | Reject a request (owner) |
| `PUT`    | `/api/exchanges/{id}/complete` | ✔ | Mark as completed (owner) — credits are transferred |
| `PUT`    | `/api/exchanges/{id}/cancel`   | ✔ | Cancel (requester or owner) — held credits are refunded |
| `POST`   | `/api/exchanges/{id}/review`   | ✔ | Review a completed exchange |

Authentication is done through the `X-User-ID: {id}` header.
Service filtering and search are performed **server-side** (query params).

### Exchange lifecycle

```
pending  ->  accepted  ->  completed
   |             |
rejected     cancelled
```

Credits are never stored as a column: the balance is the sum of the
`credit_transactions` ledger. An exchange writes to that ledger at three moments:

| Transition | Ledger entry |
|------------|--------------|
| `accepted` | `spend` — the service cost is debited from the requester (held) |
| `completed`| `earn` — the same amount is credited to the owner |
| `cancelled` from `accepted` | `refund` — the held credits go back to the requester |

`rejected`, and `cancelled` from `pending`, write nothing: no credits were held yet.

Business rules enforced by the use cases: you cannot request an exchange on your own
service (`400`), a service can only have one `pending` or `accepted` exchange at a time
(`409`), and the requester must have enough credits both to request and at acceptance
time (`400`).

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
# 409 -> {"error":"service is referenced by an exchange and cannot be deleted"}
```

An ad keeps its history: as soon as an exchange references it — even a `completed` or
`rejected` one — deletion is refused with `409`. Deactivate it instead, with
`PUT /api/services/{id}` and `{"active":false}`: no new exchange can be requested on it
(`400`), and it stops counting in `active_services`, without breaking the exchanges that
point to it. An inactive ad still appears in `GET /api/services`, flagged `"active":false`.

### A full exchange, end to end

Thierry (`id: 1`) offers "Cours de Go" (`id: 1`, 2 credits). Marie (`id: 2`) requests it,
Thierry accepts and completes it, then Marie reviews him. Both start at 10 credits.

This scenario assumes ad `1` still exists: replay it from a fresh database (`make reset`)
rather than after the delete snippet above.

**1. Marie requests the exchange** — status `pending`, no credits moved yet:

```bash
curl -i -X POST http://localhost:8000/api/exchanges \
  -H 'Content-Type: application/json' \
  -H 'X-User-ID: 2' \
  -d '{"service_id":1}'
# 201 -> {"id":1,"service_id":1,"requester_id":2,"owner_id":1,"status":"pending",...}
# 400 -> {"error":"cannot request an exchange on your own service"}   (Thierry requesting)
# 400 -> {"error":"create exchange: insufficient credits for this exchange"}
# 409 -> {"error":"create exchange: service already has an active exchange"}
```

**2. Thierry accepts** — status `accepted`, the 2 credits are held (Marie: 10 → 8):

```bash
curl -i -X PUT http://localhost:8000/api/exchanges/1/accept -H 'X-User-ID: 1'
# 200 -> {"id":1,...,"status":"accepted","updated_at":"..."}
# 403 -> {"error":"action not allowed"}                     (Marie accepting her own request)
# 400 -> {"error":"invalid exchange status transition"}     (already accepted)
```

**3. Thierry completes it** — status `completed`, the credits reach him (Thierry: 10 → 12):

```bash
curl -i -X PUT http://localhost:8000/api/exchanges/1/complete -H 'X-User-ID: 1'
# 200 -> {"id":1,...,"status":"completed"}
# 400 -> {"error":"invalid exchange status transition"}     (still pending)
```

Instead of completing, either party could cancel with
`PUT /api/exchanges/1/cancel` — Marie's 2 held credits would be refunded and her balance
back to 10.

**4. Marie reviews Thierry** — only once, only on a completed exchange:

```bash
curl -i -X POST http://localhost:8000/api/exchanges/1/review \
  -H 'Content-Type: application/json' \
  -H 'X-User-ID: 2' \
  -d '{"rating":5,"comment":"Très clair, je recommande"}'
# 201 -> {"id":1,"exchange_id":1,"author_id":2,"target_id":1,"rating":5,...}
# 400 -> {"error":"exchange must be completed before reviewing"}
# 400 -> {"error":"review already submitted for this exchange"}
# 400 -> {"error":"rating must be between 1 and 5"}
# 403 -> {"error":"only exchange participants can submit a review"}
```

**5. Thierry's dashboard reflects the whole flow**:

```bash
curl -i http://localhost:8000/api/users/1/stats
# 200 -> {"user_id":1,"active_services":1,"completed_exchanges":1,"credit_balance":12,
#         "average_rating":5,"review_count":1,"total_earned":12,"total_spent":0}
```

Marie's exchanges can be listed and filtered by status:

```bash
curl -i 'http://localhost:8000/api/exchanges?status=completed' -H 'X-User-ID: 2'
# 200 -> [{"id":1,"service_id":1,"requester_id":2,"owner_id":1,"status":"completed",...}]
# 400 -> {"error":"invalid exchange status filter"}
```

## Tests

```bash
make test
make cover
```
