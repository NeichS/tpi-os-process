PARAM ?=procesos.csv

build: 
	@echo "Compilando..."
	@go build -o bin/main cmd/main/main.go

run: build
	@./bin/main ${PARAM}	