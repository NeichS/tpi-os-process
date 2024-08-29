PARAM ?= default

build: 
	@echo "Compilando..."
	@go build -o bin/main main/main.go

run: build
	@./bin/main ${PARAM}	