run_dev:
	docker-compose -f docker-compose.yml --env-file ./docker.env.dev up -d

logs_dev:
	docker attach taskmanager

logs:
	docker attach taskmanager

kill_dev:
	docker-compose -f docker-compose.yml --env-file ./docker.env.dev down
	docker-compose -f docker-compose.yml --env-file ./docker.env.dev rm
	docker rmi taskmanager

lint:
	cd cmd/taskmanager
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && golangci-lint run

run_tests:
	go test ./cmd/taskmanager -coverprofile=coverage.out && go tool cover -html=coverage.out

