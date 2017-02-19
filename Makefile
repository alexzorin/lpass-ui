.PHONY: clean all run test

all: clean lpass-ui

clean:
	rm -f lpass-ui

lpass-ui:
	go build -o lpass-ui main.go

run:
	go run main.go

test:
	go test -v ./lpass
