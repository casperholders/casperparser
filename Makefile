#!make
include .env
export $(shell sed 's/=.*//' .env)

install_deps:
	go get -d -v .

install_pkger:
	go install github.com/markbates/pkger/cmd/pkger

pkger:
	pkger --include /sql -o cmd

binary:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/casperParser .

build: install_deps install_pkger pkger binary

test:
	go test ./... -race -covermode=atomic -coverprofile=coverage.out

codecov:
	go test ./... -race -coverprofile=coverage.out -covermode=atomic -json > report.json

report:
	go tool cover -html=coverage.out

install_migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate_db:
	migrate -source file://sql/ -database ${CASPER_PARSER_DATABASE} up

run_migrate: install_migrate migrate_db

coverage: install_deps test report