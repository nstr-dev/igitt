clean:
ifeq ($(OS),Windows_NT)
	powershell -Command "if (Test-Path out) { Remove-Item -Recurse -Force out }"
else
	rm -rf out
endif

build: clean
	go build -o out/igitt ./cmd/igitt

run: build
	./out/igitt i