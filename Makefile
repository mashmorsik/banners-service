tests:
	go test -v -race -vet=all -count=1 ./... && go test e2e_test.go

build:
	docker compose up -d --force-recreate

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2 && \
	golangci-lint run --show-stats -j=4
