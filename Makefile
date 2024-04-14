tests:
	go test -v -race -vet=all -count=1 ./...

build:
	docker compose up -d --force-recreate

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2 && \
	golangci-lint run --timeout=15s --tests --show-stats -j=4
