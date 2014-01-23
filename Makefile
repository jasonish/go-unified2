all:
	go build

test:
	go test

clean:
	go clean
	find . -name \*~ -exec rm -f {} \;
