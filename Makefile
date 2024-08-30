docker:
	docker compose up -d
	
build:
	cd server && go build -o ../bin/interview .

run: build
	./bin/interview