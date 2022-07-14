make: clean
	@go build -o gorm main.go

build:
	@go build -o gorm main.go

run:
	@go run main.go

clean:
	rm gorm
