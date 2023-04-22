APP=./bin/shortenLink

test:
	go test  ./...

build:
	go build -o $(APP) ./cmd/main.go

run_im:
	$(APP)

run_db:
	$(APP) -db

fmt:
	go fmt ./...
	goimports -l ./
	go mod tidy
