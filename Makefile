build:
	go build

clean:
	rm webcrawler.exe
	rm webcrawler

test:
	go test ./... -v -timeout 30s -cover

test-coverage:
	go test ./... -timeout 30s -coverprofile=coverage.out
	go tool cover -html=coverage.out