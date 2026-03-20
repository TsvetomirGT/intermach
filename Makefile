build:
	go build -o intermach .

run:
	go run .

install: build
	cp intermach /usr/local/bin/intermach

clean:
	rm -f intermach

.PHONY: build run install clean
