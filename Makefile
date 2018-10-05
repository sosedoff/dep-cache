build:
	go build

release:
	GOOS=linux GOARCH=amd64 go build -o ./bin/dep-cache-linux-amd64