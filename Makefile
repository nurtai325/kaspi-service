src = $(shell find . -name '*.go')

server: $(src)
	go mod tidy
	go build -o server cmd/server/main.go

docker:
	sudo docker compose up -d

test: $(src)
	go mod tidy
	go build -o test cmd/test/main.go

.PHONY: clean
clean:
	rm server test
