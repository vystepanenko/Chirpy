## Chirpy is twitter in scale. It's a guided projects from boot.dev

### Requirements
 - Go ver 1.23.2
 - Postgres installed localy or via docker

### Instalation
 - clone this repo
 - create DB for this project in Postgres
 - `cp .env.example` .env
 - fill  fields in `.env`

### Usage
```
make build_and_run
```

### Routes
	GET /api/healthz
	GET /admin/metrics
	POST /admin/reset
	POST /api/chirps
	GET /api/chirps
	GET /api/chirps/{chirpId}
	DELETE /api/chirps/{chirpId}
	POST /api/users
	POST /api/login
	POST /api/refresh
	POST /api/revoke
	PUT /api/users
	POST /api/polka/webhooks


