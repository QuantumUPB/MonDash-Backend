# MonDash Backend

<p float="left">
    <img src="upb.png" alt="University Politehnica of Bucharest" width="50"/>
    <img src="Logo.png" alt="Quantum Team @ UPB" width="100"/>
</p>

This repository contains a minimal Go backend using the [chi](https://github.com/go-chi/chi) router.

## Structure

- `cmd/` - application entry point.
- `api/` - HTTP handlers and DTOs.
- `domain/` - core domain data structures.
- `repository/` - repository interfaces.
- `services/` - service layer.
- `middlewares/` - HTTP middlewares.
- `roles.yaml` - mapping of user roles to permissions.

## Building

```bash
go build ./cmd
```

## Running

```bash
go run ./cmd
```

The server listens on `:8080` by default (or on the port defined by the `PORT`
environment variable) and exposes the following endpoints. When running via
Docker Compose the Go application listens on port `8081` and is proxied by
nginx on `http://localhost:8080`:

- `GET /healthcheck`
- `POST /update-node` - expects `{"nodes":[{"name":"<node>","status":"up|down","stored_key_count":0,"current_key_rate":0.0}]}`
- `POST /update-app`
- `POST /api/login`
- `POST /api/register` - expects `{"username":"<name>","email":"<email>","password":"<pass>","role":"<role>"}`

All non-`/api` endpoints (e.g. `/update-node`) require an `X-Auth-Token` header using the Bearer scheme, such as `X-Auth-Token: Bearer abc`.
Routes under `/api` instead rely on a cookie set by the `/api/login` endpoint. After a successful login the server returns an `auth_token` cookie that must accompany further `/api/*` requests. The in-memory authentication backend provides a default account (`admin`/`admin`) that can be used to obtain this cookie. When using MongoDB this administrator account is automatically created if the `auth_users` collection is empty.
Each endpoint currently contains placeholder logic that can be expanded later.

## Docker

The project includes a `Dockerfile` and `docker-compose.yml` for containerized
development. The Compose file starts both the backend and a MongoDB instance on
the host network. To run everything in containers:

1. Create an environment file from the template:

   ```bash
   cp .env.template .env
   ```
   The `.env` file contains a `LOG_LEVEL` variable that controls log verbosity (`debug`, `info`, `warn`, `error`).
   To send notification emails on alerts, set `EMAIL_ON_ALERT=true` and
   configure the `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD` and
   `SMTP_FROM` variables.

2. Build and start the services using Docker Compose:

   ```bash
   docker-compose up --build
   ```

   To use an existing TLS certificate instead of the self-signed one, run:

   ```bash
   make compose-up-cert CERT=/path/to/cert.crt KEY=/path/to/cert.key
   ```

The API will then be available at `http://localhost:8080` and MongoDB will be
exposed on `localhost:27017`.

## Seeding the database

The `script/populate` program fills MongoDB with the default data used by
the in-memory repositories. A separate Compose file is provided to run this
script in a container:

```bash
docker-compose -f docker-compose.populate.yml run --rm populate_db
```

This command spins up a temporary Go container that seeds the `mongodb` service
defined in the same Compose file.

## Cleaning the database

To remove all collections from the MongoDB instance you can run:

```bash
docker-compose -f docker-compose.populate.yml run --rm cleanup_db
```

This launches a short-lived container that connects to the database and drops
all data.

## Backing up the database

To create a snapshot of the MongoDB volume and a portable dump that can be
restored elsewhere, run:

```
make backup-db
```

The command stores `mongodb-data.tar.gz` and `mongodb.dump.gz` inside a new
`backup/` directory.

To import the dumped data back into the running MongoDB container, run:

```
make import-db
```

This command feeds `backup/mongodb.dump.gz` to `mongorestore`.

### Role definitions

User roles and their associated permissions are defined in `roles.yaml`. The
file maps each role to the actions it is allowed to perform. For example:

```yaml
roles:
  admin:
    - "*"
  auditor:
    - "*"
  technician:
    - "view_devices"
    - "view_nodes"
```

The `auditor` role is granted the same wildcard permission as `admin`,
providing full visibility across the system.

`LoadRolesFromEnv` in `config/roles.go` can be used to read this file, defaulting
to `roles.yaml` when the `ROLES_FILE` environment variable is unset.

# Copyright and license

This work has been implemented by Bogdan-Calin Ciobanu and Alin-Bogdan Popa under the supervision of prof. Pantelimon George Popescu, within the Quantum Team in the Computer Science and Engineering department,Faculty of Automatic Control and Computers, National University of Science and Technology POLITEHNICA Bucharest (C) 2024. In any type of usage of this code or released software, this notice shall be preserved without any changes.

If you use this software for research purposes, please follow the instructions in the "Cite this repository" option from the side panel.

This work has been partly supported by RoNaQCI, part of EuroQCI, DIGITAL-2021-QCI-01-DEPLOY-NATIONAL, 101091562.
