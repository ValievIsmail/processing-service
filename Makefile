go: 
	go run cmd/*.go
	
build:
	go build -o ./processing-service cmd/*.go

up-db:
	docker-compose up -d

down-db:
	docker-compose down

docker-build:
	docker build -t processing-service .

docker-run:
	docker run -d -p 8080:8080 processing-service