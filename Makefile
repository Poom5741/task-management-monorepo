.PHONY: dev backend frontend install test clean

dev:
	bun run dev

backend:
	cd backend && make run

frontend:
	bun run --cwd frontend dev

install:
	bun install
	cd backend && go mod download

test:
	cd backend && make test
	bun run --cwd frontend test

test-coverage:
	cd backend && make test-coverage
	bun run --cwd frontend test:coverage

build:
	cd backend && make build
	bun run --cwd frontend build

clean:
	cd backend && make clean
	rm -rf frontend/.next frontend/node_modules node_modules

lint:
	cd backend && make lint
	bun run --cwd frontend lint

fmt:
	cd backend && make fmt