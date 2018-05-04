-include .env

proto:
	@echo "==> Making proto ..."
	protoc proto/*.proto --go_out=plugins=grpc:.

start: install
	go-congress serve

install:
	go install -v ./...

install-go-deps:
	go get -u github.com/PuerkitoBio/goquery
	go get -u github.com/araddon/dateparse
	go get -u github.com/spf13/cobra

install-deps: install-migrate install-go-deps

install-migrate:
	@echo "==> Installing mattes/migrate ..."
	go get -u -d github.com/mattes/migrate/cli github.com/lib/pq
	go build -tags 'postgres' -o /usr/local/bin/migrate github.com/mattes/migrate/cli

new: reset start

reset: reset-db
	@echo "==> Resetting ..."

reset-db:
	@echo "==> Resetting database ..."
	-dropdb $(DB_NAME)
	createdb $(DB_NAME)
	#-migrate -path _migrations -database $(POSTGRES_URL) down
	migrate -path _migrations -database $(POSTGRES_URL) up

test:
	go test -v ./...

.PHONY: install test create-db proto
