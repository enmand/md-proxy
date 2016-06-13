all: clean
	go build

clean:
	-rm md-proxy


release: clean
	env GOOS=linux GOARCH=amd64 go build
