.PHONY: docker
docker: docker-clean docker-up

.PHONY: docker-clean
docker-clean:
	docker-compose down
	docker-compose rm -f

.PHONY: docker-up
docker-up:
	docker-compose up --build -d

.PHONY: test
test: docker
	go test -v -count 1 -p 1 ./...

