build:
	go build -o whistleblower .

run:
	go run .

dev:
	go run . -race

deps:
	go mod download
	go mod tidy

test:
	go test ./...

clean:
	rm -f whistleblower
	rm -f whistleblower.db

docker-build:
	docker build -t whistleblower .

docker-run:
	docker run -p 8080:8080 --env-file .env whistleblower

init-db:
	sqlite3 whistleblower.db < database/schema.sql

.PHONY: build run dev deps test clean docker-build docker-run init-db