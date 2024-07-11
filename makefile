
run:
	docker-compose -f ./docker/local/docker-compose.yml down --remove-orphans
	docker-compose -f ./docker/local/docker-compose.yml up --build	

build:
	go build -o app.exe -v ./cmd/main.go

swagger:
	go get -u github.com/swaggo/swag/cmd/swag
	swag init --parseDependency --parseInternal -g ./cmd/main.go -o docs

test:
	go test -v ./...