run-index:
	go run cmd/main.go --mode=index --index=index.json --dir=./data

run-query:
	go run cmd/main.go --mode=query --index=index.json

test:
	go test ./... -v
