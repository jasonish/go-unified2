.PHONY:	test

all:
	go build
	cd examples && go build u2bench.go
	cd examples && go build u2extract.go
	cd examples && go build simple-reader.go

test:
	go test

clean:
	go clean
	find . -name \*~ -exec rm -f {} \;
	rm -f examples/u2bench
	rm -f examples/u2extract
	rm -f examples/simple-reader
