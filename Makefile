tests:
	go test -v -race -vet=all -count=1 ./...

build:
	docker compose up -d --force-recreate