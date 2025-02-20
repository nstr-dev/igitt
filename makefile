build:
	go build -o out/igitt ./cmd/igitt

run: build
	./out/igitt i