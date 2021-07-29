build:
	go fmt ./...
	goimports -l -w `find ./ -name '*.go'`
	go build
