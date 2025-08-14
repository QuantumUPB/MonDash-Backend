.PHONY: build run docker-up seed cleanup-db backup-db import-db

build:
	go build ./cmd

run:
	go run ./cmd

compose-up:
	@[ -f .env ] || cp .env.template .env
	docker compose up --build -d

compose-up-cert:
	@[ -f .env ] || cp .env.template .env
	@[ -n "$(CERT)" ] && [ -n "$(KEY)" ] || (echo "Usage: make compose-up-cert CERT=/path/to/cert.crt KEY=/path/to/cert.key" && exit 1)
	mkdir -p docker/ssl
	cp $(CERT) docker/ssl/selfsigned.crt
	cp $(KEY) docker/ssl/selfsigned.key
	docker compose build --build-arg USE_EXISTING_CERT=true
	docker compose up -d

populate-db:
	docker compose -f docker-compose.populate.yml run --rm populate_db

cleanup-db:
	docker compose -f docker-compose.populate.yml run --rm cleanup_db

backup-db:
	mkdir -p backup
	docker compose run -T --rm -v $(PWD)/backup:/backup --entrypoint bash mongodb -c "tar czf /backup/mongodb-data.tar.gz -C /data/db ."
	docker compose exec -T mongodb mongodump --archive --gzip > backup/mongodb.dump.gz

import-db:
	@[ -f backup/mongodb.dump.gz ] || (echo "backup/mongodb.dump.gz not found. Run make backup-db first." && exit 1)
	docker compose exec -T mongodb mongorestore --archive --gzip < backup/mongodb.dump.gz

zip:
	zip -r mondash-backend.zip ./*
