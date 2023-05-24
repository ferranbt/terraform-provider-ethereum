default: install

generate:
	go generate ./...

install:
	go install .

test:
	go test -count=1 -parallel=4 ./...

testacc:
	TF_ACC=1 go test -count=1 -parallel=4 -timeout 10m -v ./...

devnet-up:
	docker-compose --project-directory ./docker_compose up -d 

devnet-down:
	docker-compose --project-directory ./docker_compose down
