clean:
ifeq ($(OS),Windows_NT)
	powershell -Command "if (Test-Path out) { Remove-Item -Recurse -Force out }"
else
	rm -rf out
endif

build: clean
ifeq ($(OS),Windows_NT)
	go build -o out/igitt.exe ./cmd/igitt
else
	go build -o out/igitt ./cmd/igitt
endif

run: build
ifeq ($(OS),Windows_NT)
	./out/igitt.exe $(ARGS)
else
	./out/igitt $(ARGS)
endif